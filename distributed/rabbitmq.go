package distributed

import (
	"fmt"
	"jzmall/common"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// name of queues
	QueueName string
	// exchange machine
	Exchange string
	// key
	Key string
	// connection info
	Mqurl string
	sync.Mutex
}

type IConsumeHandler interface {
	DoConsume(dl *amqp.Delivery)
}

// create a new instance of RabbitMQ client
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	amqpAuthenInfo, err := common.ReadPrivateFile("distributed/AMQP")
	if err != nil {
		panic(err)
	}
	amqpUrl := "amqp://" + amqpAuthenInfo[0] + ":" + amqpAuthenInfo[1] + "@127.0.0.1:5672/" + amqpAuthenInfo[2]
	rabbitmq := RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: amqpUrl}
	// create a connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "Failed to connect!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "Failed to get channel")
	return &rabbitmq
}

// disconnect channel and connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// error handling
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// simple pattern: 1. create object of RabbitMQ under simple pattern
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

// simple pattern: 2. publish code
func (r *RabbitMQ) PublishSimple(message string) error {
	// 1. request for a queue, skip if it doesn't exist, otherwise create a new queue
	r.Lock()
	defer r.Unlock()
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // extra
	)
	if err != nil {
		return err
	}
	// 2. send message to the queue
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, // mandatory, if no valid queue found, send back to producer
		false, // immediate, if no consumer found for the queue, send back to producer
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return nil
}

// simple pattern: 3. Consume code
func (r *RabbitMQ) ConsumeSimple(consumeHandle IConsumeHandler) {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // extra
	)
	if err != nil {
		fmt.Println(err)
	}

	r.channel.Qos(
		1,     // maximum number of messages consumed each time
		0,     // maximum size of messages consumed each time
		false, // only apply to this channel
	)

	// receive messages
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",    // used for identifying different consumers
		false, // autoAck
		false, // exclusive
		false, // noLocal, if true, we can't pass messages from producers and consumers in the same connection
		false, // noWait
		nil,
	)

	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// start goroutine to handle messages
	go func() {
		for d := range msgs {
			// the main logics used to deal with messages
			consumeHandle.DoConsume(&d)
		}
	}()

	log.Println("[*] Waiting for messages")
	<-forever
}

// publish/subscribe pattern: 1. create an instance
func NewRabbitMQPubSub(exchange string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchange, "")
	var err error
	// create a connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "Failed to connect!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "Failed to get channel")
	return rabbitmq
}

func (r *RabbitMQ) PublishSub(message string) {
	// publish/subscribe pattern: 2. declare an exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", // kind
		true,     // durable
		false,
		false, // internal, whether the exchange should only be used for binding other exchanges
		false, // noWait
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")

	//	3. send the message
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		fmt.Println(err)
	}
}

// publish/subscribe pattern: 4. consume code
func (r *RabbitMQ) ConsumePubSub() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")

	// try to create a queue
	q, err := r.channel.QueueDeclare(
		"", // empty, cuz rabbitmq will assign a random name to it
		false,
		false,
		true,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare a queue")

	// bind the queue to the exchange
	err = r.channel.QueueBind(
		q.Name,
		"", // for pub/sub pattern, key must be empty!
		r.Exchange,
		false, // noWait
		nil,
	)

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", string(d.Body))
		}
	}()

	fmt.Println("Listening to RabbitMQ ... ")
	<-forever
}

func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// create a connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "Failed to connect!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "Failed to get channel")
	return rabbitmq
}

func (r *RabbitMQ) PublishRouting(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", // kind
		true,     // durable
		false,
		false, // internal, whether the exchange should only be used for binding other exchanges
		false, // noWait
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")

	//	3. send the message
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		fmt.Println(err)
	}
}

func (r *RabbitMQ) ConsumeRouting() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")

	// try to create a queue
	q, err := r.channel.QueueDeclare(
		"", // empty, cuz rabbitmq will assign a random name to it
		false,
		false,
		true,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare a queue")

	// bind the queue to the exchange
	err = r.channel.QueueBind(
		q.Name,
		r.Key, // for pub/sub pattern, key must be empty!
		r.Exchange,
		false, // noWait
		nil,
	)

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", string(d.Body))
		}
	}()

	fmt.Println("Listening to RabbitMQ ... ")
	<-forever
}

func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// create a connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "Failed to connect!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "Failed to get channel")
	return rabbitmq
}

func (r *RabbitMQ) PublishTopic(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic", // kind
		true,    // durable
		false,
		false, // internal, whether the exchange should only be used for binding other exchanges
		false, // noWait
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")

	//	3. send the message
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		fmt.Println(err)
	}
}

func (r *RabbitMQ) ConsumeTopic() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")

	// try to create a queue
	q, err := r.channel.QueueDeclare(
		"", // empty, cuz rabbitmq will assign a random name to it
		false,
		false,
		true,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare a queue")

	// bind the queue to the exchange
	// the binding key match rule is :
	// words are separated by .
	// * can match one word, # can match multiple words (0 words applies)
	err = r.channel.QueueBind(
		q.Name,
		r.Key, // for pub/sub pattern, key must be empty!
		r.Exchange,
		false, // noWait
		nil,
	)

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", string(d.Body))
		}
	}()

	fmt.Println("Listening to RabbitMQ ... ")
	<-forever
}
