package main

import (
	"apps/services"
	"fmt"
	"log"
	"os"
	"sync"
)

type Env struct {
    Rmq *services.Rmq
    Hetzner *services.Hetzner
}

func main() {


    rmq := services.RabbitmqSetup()
    hetzner := services.HetznerSetup();
    env := &Env{
        Rmq: rmq,
        Hetzner: hetzner,
    }

    joinCommand := os.Getenv("JOIN_COMMAND")

    // Do k8s stuff
    for {
        currentNumberOfNodes := services.NumberOfNodes()
        lengthOfQueue := env.Rmq.QueueLength()
    
        numberOfNodesToMake := (lengthOfQueue / 100) - currentNumberOfNodes
        if numberOfNodesToMake > 20 {
            log.Panicf("Cannot create %d!", numberOfNodesToMake)
        }

        var wg sync.WaitGroup
  
        if numberOfNodesToMake > 0 {
            wg.Add(numberOfNodesToMake)
            for  i := 0; i < numberOfNodesToMake; i++ {
                go func() {
                    defer wg.Done()
                    env.Hetzner.CreateNode(joinCommand)
                }()
            }
        }
       
        // Wait till all new nodes have initalised
        wg.Wait()
    }

}
