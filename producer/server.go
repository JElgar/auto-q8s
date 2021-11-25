package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    
    "github.com/streadway/amqp"
)

type Env struct {
    queue amqp.Queue
    channel *amqp.Channel
}

func (env *Env) handlePostRequest(w http.ResponseWriter, r *http.Request){
    body, err := ioutil.ReadAll(r.Body)
    if (err != nil) {
        fmt.Fprintf(w, "Failed to read request", err)
        log.Fatal(http.ListenAndServe("Error processing request", nil))
        return
    }
    env.publishToQueue(body)
}

func (env *Env) handleRequests() {
    http.HandleFunc("/", env.handlePostRequest)
    log.Fatal(http.ListenAndServe(":10000", nil))
}

func rabbitmqSetup() (*amqp.Channel, amqp.Queue) {
    conn, rabbitErr := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if (rabbitErr != nil) {
        log.Fatal("Failed to connect to rabbitmq")
    }

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

    return channel, queue
}

func (env *Env) publishToQueue(data []byte) {
    err := env.channel.Publish(
        "",     // exchange
        env.queue.Name, // routing key
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

func main() {
    channel, queue := rabbitmqSetup()
    env := &Env{queue: queue, channel: channel}

    env.handleRequests()
}
