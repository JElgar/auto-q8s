package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "encoding/json"
    
    "github.com/streadway/amqp"
)

type Request struct {
    Id          string `json:"id"`
    Location    string `json:"location"`
    Amount      string `json:"amount"`
}

conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

func pushToQueue(w http.ResponseWriter, r *http.Request){
    body, err := ioutil.ReadAll(r.Body)
    if (err != nil) {
        fmt.Fprintf(w, "Failed to read request", err)
        log.Fatal(http.ListenAndServe("Error processing request", nil))
        return
    }
   
    var request Request 
    err = json.Unmarshal(body, &request)
    if (err != nil) {
        log.Fatal(http.ListenAndServe("Error parsing request", nil))
        fmt.Fprintf(w, "Failed to parse request", err)
        return
    }

    fmt.Println("Endpoint Hit: Id: ")
    fmt.Println(request.Id)
}

func handleRequests() {
    http.HandleFunc("/", pushToQueue)
    log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
    handleRequests()
}
