package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"os"
)

func main() {
	// Read a source file and place it on a rabbit topic exchange
	// Exchange , Queue and binding must be in place before use

	var url, exchange, routingKey, encryptionKey, signingKey, name, password string
	var port, host, queue, vhost, messageFilePath string
	var msgBody []byte
	var err error
	const PRINT_MSG_CHAR_COUNT = 20 // only print the first few characters of the message

	// access command line parameters

	flag.StringVar(&name, "n", "", "name of the rabbit user")
	flag.StringVar(&password, "p", "", "password of the rabbit user")
	flag.StringVar(&port, "rport", "", "port used to connect to rabbit")
	flag.StringVar(&host, "rhost", "", "hostname used to connect to rabbit")
	flag.StringVar(&vhost, "v", "", "vhostname used to connect to rabbit")
	flag.StringVar(&url, "u", "", "url connection string ")
	flag.StringVar(&queue, "q", "", "name of the rabbit queue")
	flag.StringVar(&exchange, "x", "", "name of the rabbit exchange")
	flag.StringVar(&routingKey, "r", "", "rabbit routing key")
	flag.StringVar(&encryptionKey, "e", "", "path to a private key file used for encryption")
	flag.StringVar(&signingKey, "s", "", "path to a private key used for signing")
	flag.StringVar(&messageFilePath, "f", "", "path to filename to send")

	flag.Parse()

	// If no value given on command line , then look at environment variables , else use a default value

	name = getFromEnvIfEmpty(name, "RABBITMQ_DEFAULT_USER", "guest")
	password = getFromEnvIfEmpty(password, "RABBITMQ_DEFAULT_PASS", "guest")
	port = getFromEnvIfEmpty(port, "RABBITMQ_PORT", "5672")
	host = getFromEnvIfEmpty(host, "RABBITMQ_HOST", "localhost")
	queue = getFromEnvIfEmpty(queue, "RABBIT_SURVEY_QUEUE", "rabbit")
	exchange = getFromEnvIfEmpty(exchange, "RABBITMQ_EXCHANGE", "message")
	vhost = getFromEnvIfEmpty(vhost, "RABBITMQ_DEFAULT_VHOST", "%2f")

	if url == "" {
		url = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", name, password, host, port, vhost)
	} else {
		fmt.Println("url use overrides specific parameters") //A specific URL is set , so use it
	}

	msgBody, err = getBody(messageFilePath)
	failOnError(err, "could not read message body")

	// If encyrpt specified then encrypt
	// if sign specified then sign

	err = sendToRabbit(url, exchange, queue, routingKey, msgBody)
	failOnError(err, "unable to send message to rabbitmq")

	var msgSize = len(msgBody)
	if msgSize < PRINT_MSG_CHAR_COUNT {
		fmt.Println(fmt.Sprintf("message:'%s' (len=%d) published to exchange:'%s' using routing key:'%s'", string(msgBody), msgSize, exchange, routingKey))
	} else {
		fmt.Println(fmt.Sprintf("message:'%s...' (len=%d) published to exchange:'%s' using routing key:'%s'", string(msgBody[0:PRINT_MSG_CHAR_COUNT]), msgSize, exchange, routingKey))
	}
}

func getFromEnvIfEmpty(target string, key string, defaultValue string) string {

	if target != "" {
		return target
	}

	if value, present := os.LookupEnv(key); present {
		return value
	}

	return defaultValue
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

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("%s: %s", msg, err)
		os.Exit(1) // No specific err codes atm
	}
}

func sendToRabbit(url string, exchange string, queue string, routingKey string, msgBody []byte) error {

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
