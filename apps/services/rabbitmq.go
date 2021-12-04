package services

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type Rmq struct {
    Queue amqp.Queue
    Channel *amqp.Channel
}

func (rmq *Rmq) Publish(data []byte) {
    err := rmq.Channel.Publish(
        "",     // exchange
        rmq.Queue.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing {
        ContentType: "text/plain",
        Body:        data,
    })
    
    if (err != nil) {
        log.Fatal("Failed to publish to queue")
    }
}


func (rmq *Rmq) Consumer() (<-chan amqp.Delivery) {
    data, err := rmq.Channel.Consume(
      rmq.Queue.Name, // queue
      "",     // consumer
      true,   // auto-ack
      false,  // exclusive
      false,  // no-local
      false,  // no-wait
      nil,    // args
    )
    
    if (err != nil) {
        log.Fatal("Failed to publish to queue")
    }

    return data
}


func (rmq *Rmq) Consume() []byte {
    data, err := rmq.Channel.Consume(
      rmq.Queue.Name, // queue
      "",     // consumer
      true,   // auto-ack
      false,  // exclusive
      false,  // no-local
      false,  // no-wait
      nil,    // args
    )
    
    if (err != nil) {
        log.Fatal("Failed to publish to queue")
    }

    return (<-data).Body
}


func connectionLoop(connectionStr string) *amqp.Connection {
    for {
      conn, rabbitErr := amqp.Dial(connectionStr)
      if (rabbitErr == nil) {
        log.Println("Connected to rabbitmq!")
        return conn
      }
      log.Println("Failed to connect to cluster. Trying again in 2 seconds.")
      time.Sleep(2 * time.Second)
    }
}


func RabbitmqSetup() *Rmq {
    connectionStr := fmt.Sprintf(
        "amqp://%s:%s@%s:%s/",
        os.Getenv("RMQ_USER"),
        os.Getenv("RMQ_PASSWORD"),
        os.Getenv("RMQ_HOST"),
        os.Getenv("RMQ_PORT"),
    )
    conn := connectionLoop(connectionStr)

    channel, err := conn.Channel()
    if (err != nil) {
        log.Fatal("Failed to create channel")
    }

    queue, err := channel.QueueDeclare(
      "hello", // name
      false,   // durable
      false,   // delete when unused
      false,   // exclusive
      false,   // no-wait
      nil,     // arguments
    )
    
    if (err != nil) {
        log.Fatal("Failed to create queue")
    }

    return &Rmq{Queue: queue, Channel: channel} 
}

