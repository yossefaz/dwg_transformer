package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type rabbitmq struct {
	conn  *amqp.Connection
	chanL *amqp.Channel
	queue amqp.Queue
}

func handleError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}

}

type pickFile struct {
	name string
	path string
}

func newRabbit(connString string, queueName string) (instance rabbitmq) {
	conn, err := amqp.Dial(connString)
	handleError(err, "Can't connect to AMQP")
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	queue, err := amqpChannel.QueueDeclare(queueName, true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")
	return rabbitmq{
		conn:  conn,
		chanL: amqpChannel,
		queue: queue,
	}
}

func (rmq rabbitmq) listenMessage() {

	err := rmq.chanL.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")
	messageChannel, err := rmq.chanL.Consume(
		rmq.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		counter := 0
		fmt.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			counter++
			fmt.Printf("\nReceived a message: %s \n", d.Body)

			pFIle := &pickFile{}

			err := json.Unmarshal(d.Body, pFIle)

			if err != nil {
				fmt.Printf("Error decoding JSON: %s", err)
			}
			if err := d.Ack(false); err != nil {
				fmt.Printf("Error acknowledging message : %s", err)
			} else {
				fmt.Println("Acknowledged message :", counter)
			}

		}
	}()

	// Stop for program termination
	<-stopChan

}
