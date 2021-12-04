package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
    "apps/services"
    "encoding/json"
)

type Env struct {
    Rmq *services.Rmq
}

func (env *Env) handlePostRequest(w http.ResponseWriter, r *http.Request){
    body, err := ioutil.ReadAll(r.Body)
    if (err != nil) {
        fmt.Fprintf(w, "Failed to read request, %s", err)
        log.Fatal(http.ListenAndServe("Error processing request", nil))
        return
    }
    env.Rmq.Publish(body)
    
    resp := make(map[string]string)
    resp["message"] = "Published to queue"

    w.WriteHeader(http.StatusOK)
    jsonResp, _ := json.Marshal(resp)
    w.Write(jsonResp)
}

func (env *Env) handleRequests() {
    http.HandleFunc("/", env.handlePostRequest)
    log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
    rmq := services.RabbitmqSetup()
    env := &Env{Rmq: rmq}

    env.handleRequests()
}
