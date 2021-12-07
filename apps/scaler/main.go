package main

import (
	"apps/services"
	"fmt"
	// "os"
	"time"
)

type Env struct {
    Rmq *services.Rmq
    Hetzner *services.Hetzner
}

func main() {

    // Do k8s stuff
    services.DoStuff()

    rmq := services.RabbitmqSetup()
    hetzner := services.HetznerSetup();
    env := &Env{
        Rmq: rmq,
        Hetzner: hetzner,
    }

    // env.Hetzner.DeleteNode("worker-node-098c185b-5d4a-479a-556a-53f3144695ff")
    // env.Hetzner.CreateNode(os.Getenv("JOIN_COMMAND"))
    for {
        time.Sleep(time.Second * 5)
        fmt.Println("The queue length is: ")
        fmt.Println(env.Rmq.QueueLength())
    }
}
