package main
import(
		"fmt"
		"flag"
		"github.com/streadway/amqp"
		"os"
		"io/ioutil"
)

func main() {

	// Read a source file and place it on a rabbit topic exchange
	var url, exchange, routing_key,    encryption_key,    signing_key, name, password string
	var port, host , queue, vhost, message_path_and_file string
	var msg_body []byte
	var err error

	getArguments(&name, &password, &port, &host, &url, &queue, &exchange, &routing_key, &encryption_key, &signing_key, &vhost, &message_path_and_file)

	msg_body = getBody(message_path_and_file)

	// If encyrpt specified then encrypt
	// if sign specified then sign

	sendToRabbit(url, exchange, queue, routing_key, err, msg_body)
	
}

func getArguments(name *string,
	password *string,
	port *string,
	host *string,
	url *string,
	queue *string,
	exchange *string,
	routing_key *string,
	encryption_key *string,
	signing_key *string,
	vhost *string,
	message_path_and_file *string)  {
	// define the command line parameters
	flag.StringVar(name, "n", "", "name of the rabbit user")
	flag.StringVar(password, "p", "", "password of the rabbit user")
	flag.StringVar(port, "rport", "", "port used to connect to rabbit")
	flag.StringVar(host, "rhost", "localhost", "hostname used to connect to rabbit")
	flag.StringVar(vhost, "v", "", "vhostname used to connect to rabbit")
	flag.StringVar(url, "u", "", "url connection string ")
	flag.StringVar(queue, "q", "", "name of the rabbit queue")
	flag.StringVar(exchange, "x", "message", "name of the rabbit exchange")
	flag.StringVar(routing_key, "r", "sdx-made-up-routing-key", "rabbit routing key")
	flag.StringVar(encryption_key, "e", "sdx-made-encryption_key", "path to a private key file used for encryption")
	flag.StringVar(signing_key, "s", "sdx-made-signing_key", "path to a private key used for signing")
	flag.StringVar(message_path_and_file, "f", "", "path to filename to send")
	flag.Parse()
	// Use values from environment variables if no value empty
	getFromEnvIfEmpty(name, "RABBITMQ_DEFAULT_USER", "guest")
	getFromEnvIfEmpty(password, "RABBITMQ_DEFAULT_PASS", "guest")
	getFromEnvIfEmpty(port, "RABBITMQ_PORT", "5672")
	getFromEnvIfEmpty(host, "RABBITMQ_HOST", "rabbit")
	getFromEnvIfEmpty(queue, "RABBIT_SURVEY_QUEUE", "rabbit")
	getFromEnvIfEmpty(exchange, "RABBITMQ_EXCHANGE", "rabbit")
	getFromEnvIfEmpty(vhost, "RABBITMQ_DEFAULT_VHOST", "%2f")
	if *url == "" {
		*url = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", *name, *password, *host, *port, *vhost)
	} else {
		fmt.Println("Url use overrides specific parameters") //A specific URL is set , so use it
	}
}

func getFromEnvIfEmpty(target *string, key string, defaultValue string)  {
	if *target == "" {
		if value, present := os.LookupEnv(key); present {
			*target = value
		} else {
			*target = defaultValue
		}
	}
}

// Consider adding StdIn reading here to support piping ?
func getBody(file_path string)([]byte){
	var msg_body []byte
	var err error
	if len(file_path) != 0 {
		msg_body, err = ioutil.ReadFile(file_path)
		failOnError(err, "Cannot read from file")
	} else {panic("No file name supplied")}  // Only mandatory until we support piping
	return msg_body
}

func failOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func sendToRabbit(url string, exchange string, queue string, routing_key string, err error, msg_body []byte) {
	// Connect to Rabbit
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer conn.Close()
	defer ch.Close()


	// REMOVE THIS CODE !!!  Do not Declare in prod - must exist Prior to run ?
	rabbitPrepare(ch, exchange, queue, routing_key)
	// END OF CODE TO REMOVE


	err = ch.Publish(
		exchange,    // exchange
		routing_key, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body: msg_body,
		})

	failOnError(err, "Failed to publish a message")
}




// Not for Prod
func rabbitPrepare(ch *amqp.Channel, exchange string, queue string, routing_key string)  {
	err := ch.ExchangeDeclare(
		exchange, // name
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	q, err := ch.QueueDeclare(
		queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	Use(q)
	failOnError(err, "Could not create queue ")
	err = ch.QueueBind(
		queue,       // queue name
		routing_key, // routing key
		exchange,    // exchange
		false,
		nil)
	failOnError(err, fmt.Sprintf("could not bind queue %s to exchange %s", queue, exchange))
}

func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
// End of Not for Prod



