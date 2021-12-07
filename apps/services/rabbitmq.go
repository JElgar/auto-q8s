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
    log.Printf("Published to queue")
}


func (rmq *Rmq) Consumer() (<-chan amqp.Delivery) {
    data, err := rmq.Channel.Consume(
      rmq.Queue.Name, // queue
      "",     // consumer
      false,   // auto-ack
      false,  // exclusive
      false,  // no-local
      false,  // no-wait
      nil,    // args
    )
    
    if (err != nil) {
        log.Fatal("Failed to consume from queue")
    }

    return data
}


func (rmq *Rmq) Consume() []byte {
    data := rmq.Consumer()
    return (<-data).Body
}


func (rmq *Rmq) QueueLength() int {
  // rmq.Publish([]byte("test"))
  // result, ok, err := rmq.Channel.Get(rmq.Queue.Name, false)
  // if (!ok) {
  //   fmt.Println("UHOH")
  //   fmt.Println(err)
  //   fmt.Println(ok)
  //   fmt.Println(result)
  //   fmt.Println("Failed to get queue length")
  // }
  // return result.MessageCount
  rmq.refreshQueue(rmq.Queue.Name)
  return rmq.Queue.Messages
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

func (rmq *Rmq) refreshQueue(queueName string) {
  queue, err := rmq.Channel.QueueDeclarePassive(
    queueName, // name
    false,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )

  if (err != nil) {
      log.Fatal("Failed to refresh queue")
  }

  rmq.Queue = queue
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
    channel.Qos(1, 0, true)
    if (err != nil) {
        log.Fatal("Failed to create channel")
    }
    rmq := &Rmq{Channel: channel} 
  
    queue, err := channel.QueueDeclare(
      "hello", // name
      false,   // durable
      false,   // delete when unused
      false,   // exclusive
      false,   // no-wait
      nil,     // arguments
    )

    if (err != nil) {
      rmq.refreshQueue("hello")
    } else {
      rmq.Queue = queue 
    }

    return rmq
}

