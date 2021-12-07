package main

import (
	"apps/services"
	"fmt"
	"log"
	"runtime"
	"time"
)

func idle() {
    done := make(chan int)
    
    for i := 0; i < runtime.NumCPU(); i++ {
        go func() {
            for {
                select {
                case <-done:
                    return
                default:
                }
            }
        }()
    }
    
    time.Sleep(time.Second * 10)
    close(done)
}

func main() {
    rmq := services.RabbitmqSetup()
    fmt.Println("Here we go")

    for d := range rmq.Consumer() {
        log.Printf("Doing work on: %s", d.Body)
        idle()
        log.Printf("Done work on: %s", d.Body)
        rmq.Channel.Ack(d.DeliveryTag, false)
    }
}
