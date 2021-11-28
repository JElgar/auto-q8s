package main

import (
	"log"
    "apps/services"
)

func main() {
    rmq := services.RabbitmqSetup()
    for d := range rmq.Consumer() {
        log.Printf("Doing work on: %s", d.Body)
    }
}
