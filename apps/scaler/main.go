package main

import (
	"apps/services"
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
        log.Printf("Checking")
        currentNumberOfNodes := services.NumberOfNodes()
        lengthOfQueue := env.Rmq.QueueLength()
    
        numberOfNodesToMake := int((lengthOfQueue / 100) - currentNumberOfNodes)
        if numberOfNodesToMake > 20 {
            log.Panicf("Cannot create %d!", numberOfNodesToMake)
        }
        log.Printf("Number of nodes in cluster: %d", currentNumberOfNodes)
        log.Printf("Number of items in queue: %d", lengthOfQueue)
        log.Printf("Number of nodes to create: %d", numberOfNodesToMake)

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
        log.Print("Waiting for nodes to be created and inited")
        wg.Wait()
        log.Println("Done")
    }

}
