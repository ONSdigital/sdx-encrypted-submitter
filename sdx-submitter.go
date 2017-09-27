package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config part read from command line and part from sdx-submitter.yml

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

const configFileName = "./sdx-submitter.yml"

func main() {

	// Read a source file and place it on a rabbit topic exchange
	// Exchange , Queue and binding must be in place before use

	var config config

	// access command line parameters

	flag.StringVar(&config.Name, "n", "", "name of the rabbit user")
	flag.StringVar(&config.Password, "p", "", "password of the rabbit user")
	flag.StringVar(&config.EncryptionKeyFile, "e", "", "path to a private key file used for encryption")
	flag.StringVar(&config.SigningKeyFile, "s", "", "path to a private key used for signing")
	flag.StringVar(&config.MessageFilePath, "f", "", "path to filename to send")

	flag.Parse()

	// Get config file values

	configFile, filepathError := filepath.Abs(configFileName)
	exitOnError(filepathError, fmt.Sprintf(" cannot get absolute filename from %s", configFileName))

	yamlFile, readfileError := ioutil.ReadFile(configFile)
	exitOnError(readfileError, fmt.Sprintf("unable to read from %s", configFileName))

	marshalError := yaml.Unmarshal(yamlFile, &config)
	exitOnError(marshalError, fmt.Sprintf("unable to unMarshal yaml from %s", configFileName))

	message, messageError := getMessage(config.MessageFilePath)
	exitOnError(messageError, "could not read message body")

	// If encrypt specified then encrypt
	// if sign specified then sign

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Name, config.Password, config.Host, config.Port, config.Vhost)
	rabbitError := sendToRabbit(url, config.Exchange, config.RoutingKey, message)
	exitOnError(rabbitError, "unable to send message to rabbitmq")

	fmt.Printf("message from file:'%s' published to exchange:'%s' using routing key:'%s\n", config.MessageFilePath, config.Exchange, config.RoutingKey)
}

func exitOnError(err error, msg string) {
	if err != nil {
		fmt.Println("%s: %s", msg, err)
		os.Exit(1)
	}
}

// Consider adding stdin reading here to support piping ?
func getMessage(filePath string) ([]byte, error) {
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

func sendToRabbit(url string, exchange string, routingKey string, msgBody []byte) error {

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

	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgBody,
		})

	return err
}
