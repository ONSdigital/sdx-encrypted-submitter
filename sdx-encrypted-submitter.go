package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sdx-encrypted-submitter/authentication"
)

// Config part read from command line and part from sdx-encrypted-submitter.yml

type config struct {
	Name              string
	Password          string
	Port              int    `yaml:"port"`
	Host              string `yaml:"host"`
	Exchange          string `yaml:"exchange"`
	Vhost             string `yaml:"vhost"`
	RoutingKey        string `yaml:"routingkey"`
	EncryptionKeyFile string
	SigningKeyFile    string
	MessageFilePath   string
}

const configFileName = "./sdx-encrypted-submitter.yml"

var testArgs []string // testArgs not exported , used for testing only

func main() {

	// Read a source file and place it on a rabbit topic exchange
	// Exchange , Queue and binding must be in place before use

	var config config
	var txID string

	//  Use testArgs if supplied ahead of command line
	a := os.Args[1:]
	if testArgs != nil {
		a = testArgs
	}

	// access command line parameters

	flag.StringVar(&config.Name, "n", "", "name of the rabbit user")
	flag.StringVar(&config.Password, "p", "", "password of the rabbit user")
	flag.StringVar(&config.EncryptionKeyFile, "e", "", "path to a public key file used for encryption")
	flag.StringVar(&config.SigningKeyFile, "s", "", "path to a private key used for signing")
	flag.StringVar(&config.MessageFilePath, "f", "", "path to filename to send")

	flag.CommandLine.Parse(a)

	if config.EncryptionKeyFile == "" {
		exitOnError(errors.New("encryption key file not supplied"), "encryption key required")
	}

	if config.SigningKeyFile == "" {
		exitOnError(errors.New("signing key file not supplied"), "signing key required")
	}

	// Get config file values

	configFile, filepathError := filepath.Abs(configFileName)
	exitOnError(filepathError, fmt.Sprintf(" cannot get absolute filename from %s", configFileName))

	yamlFile, readfileError := ioutil.ReadFile(configFile)
	exitOnError(readfileError, fmt.Sprintf("unable to read from %s", configFileName))

	marshalError := yaml.Unmarshal(yamlFile, &config)
	exitOnError(marshalError, fmt.Sprintf("unable to unMarshal yaml from %s", configFileName))

	message, messageError := getRawMessage(config.MessageFilePath)
	exitOnError(messageError, "could not read message body")

	var mappedData map[string]interface{}
	fileErr := json.Unmarshal(message, &mappedData) //file contents are arbitrary Json
	txID = fmt.Sprintf("%v", mappedData["tx_id"])
	exitOnError(fileErr, "Could not marshal Json from input file")

	jwe, tokenError := authentication.GetJwe(mappedData, config.SigningKeyFile, config.EncryptionKeyFile)
	message = []byte(jwe)
	if tokenError != nil {
		exitOnError(errors.New(""), fmt.Sprintf("%s %s", (*tokenError).From, (*tokenError).Desc))
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Name, config.Password, config.Host, config.Port, config.Vhost)
	rabbitError := sendToRabbit(url, config.Exchange, config.RoutingKey, txID, message)
	exitOnError(rabbitError, "unable to send message to rabbitmq")

	fmt.Printf("message from file:'%s' published to exchange:'%s' using routing key:'%s\n", config.MessageFilePath, config.Exchange, config.RoutingKey)
}

func exitOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, " - ", err)
		os.Exit(1)
	}
}

//TODO Consider adding stdin reading here to support piping ?
func getRawMessage(filePath string) ([]byte, error) {
	var msgBody []byte
	var err error

	if len(filePath) == 0 {
		return nil, errors.New("no file name supplied")
	}

	msgBody, err = ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return msgBody, nil
}

func sendToRabbit(url string, exchange string, routingKey string, txID string, msgBody []byte) error {

	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	conn, err = amqp.Dial(url)
	defer conn.Close()
	if err != nil {
		return err
	}

	ch, err = conn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}

	var headers = amqp.Table{"tx_id": txID}
	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgBody,
			Headers:     headers,
		})

	return err
}
