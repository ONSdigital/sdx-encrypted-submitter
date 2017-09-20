package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"os"
	"gopkg.in/yaml.v2"
	"path/filepath"
)


// config part read from command line and part from sdx-submitter.yml

type Config struct{
	Name string
	Password string
	Port int
	Host string
	Queue string
	Exchange string
	Vhost string
	RoutingKey string
	Url string
	EncryptionKeyFile string
	SigningKeyFile string
	MessageFilePath string
}

func main() {

	// Read a source file and place it on a rabbit topic exchange
	// Exchange , Queue and binding must be in place before use

	var config Config
	var msgBody []byte
	var yamlFile []byte
	var err error
	const printMsgCharCount = 20 // only print the first few characters of the message
	const configFileName = "./sdx-submitter.yml"

	// access command line parameters
	flag.StringVar(&config.Name, "n", "", "name of the rabbit user")
	flag.StringVar(&config.Password, "p", "", "password of the rabbit user")
	flag.StringVar(&config.EncryptionKeyFile, "e", "", "path to a private key file used for encryption")
	flag.StringVar(&config.SigningKeyFile, "s", "", "path to a private key used for signing")
	flag.StringVar(&config.MessageFilePath, "f", "", "path to filename to send")

	flag.Parse()

	// Get config file values

	configFile, err := filepath.Abs(configFileName)
	exitOnError(err,fmt.Sprintf(" cannot open %s",configFileName))

	yamlFile, err = ioutil.ReadFile(configFile)
	exitOnError(err, fmt.Sprintf("unable to read from %s",configFileName))

	err = yaml.Unmarshal(yamlFile, &config)
	exitOnError(err,fmt.Sprintf("unable to unMarshal yaml from %s",configFileName))

	config.Url = fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Name, config.Password, config.Host, config.Port, config.Vhost)

	msgBody, err = getBody(config.MessageFilePath)
	exitOnError(err, "could not read message body")

	// If encyrpt specified then encrypt
	// if sign specified then sign

	err = sendToRabbit(config.Url, config.Exchange, config.RoutingKey, msgBody)
	exitOnError(err, "unable to send message to rabbitmq")

	var msgSize = len(msgBody)
	if msgSize < printMsgCharCount {
		fmt.Println(fmt.Sprintf("message:'%s' (len=%d) published to exchange:'%s' using routing key:'%s'", string(msgBody), msgSize, config.Exchange, config.RoutingKey))
	} else {
		fmt.Println(fmt.Sprintf("message:'%s...' (len=%d) published to exchange:'%s' using routing key:'%s'", string(msgBody[0:printMsgCharCount]), msgSize, config.Exchange, config.RoutingKey))
	}
}


// Consider adding stdin reading here to support piping ?
func getBody(file_path string) ([]byte, error) {
	var msgBody []byte
	var err error

	if len(file_path) == 0 {
		return nil, errors.New("no file name supplied")
	}

	msgBody, err = ioutil.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	return msgBody, err
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

func exitOnError(err error, msg string) {
	if err != nil {
		fmt.Println("%s: %s", msg, err)
		os.Exit(1)
	}
}
