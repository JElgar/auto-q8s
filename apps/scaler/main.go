package main

import (
    "apps/services"
)

type Env struct {
    Rmq *services.Rmq
    Hetzner *services.Hetzner
}

func main() {
    // rmq := services.RabbitmqSetup()
    hetzner := services.HetznerSetup();
    env := &Env{
    //    Rmq: rmq,
        Hetzner: hetzner,
    }

    // env.Hetzner.DeleteNode("worker-node-098c185b-5d4a-479a-556a-53f3144695ff")
    // env.Hetzner.CreateNode()
    env.Hetzner.Test()
}
