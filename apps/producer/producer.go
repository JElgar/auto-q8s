package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
    "apps/services"
)

type Env struct {
    Rmq *services.Rmq
}

func (env *Env) handlePostRequest(w http.ResponseWriter, r *http.Request){
    body, err := ioutil.ReadAll(r.Body)
    if (err != nil) {
        fmt.Fprintf(w, "Failed to read request", err)
        log.Fatal(http.ListenAndServe("Error processing request", nil))
        return
    }
    env.Rmq.Publish(body)
}

func (env *Env) handleRequests() {
    http.HandleFunc("/", env.handlePostRequest)
    log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
    rmq := services.RabbitmqSetup()
    env := &Env{Rmq: rmq}

    env.handleRequests()
}
