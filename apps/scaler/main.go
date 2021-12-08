package main

import (
	"apps/services"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type Env struct {
    Rmq *services.Rmq
    Hetzner *services.Hetzner
    K8s *services.K8sEnv
}

func main() {

    hetzner := services.HetznerSetup();
    rmq := services.RabbitmqSetup()
    k8sEnv := services.InitK8s()
    env := &Env{
        Rmq: rmq,
        Hetzner: hetzner,
        K8s: k8sEnv,
    }

    joinCommand := os.Getenv("JOIN_COMMAND")

    // Do k8s stuff
    for {
        log.Printf("Checking")
        currentNumberOfNodes := env.K8s.NumberOfNodes()
        lengthOfQueue := env.Rmq.QueueLength()
    
        numberOfNodesToMake := int(math.Ceil(float64(lengthOfQueue) / float64(100))) - currentNumberOfNodes
        deploymentSize := int(math.Max(math.Ceil(float64(lengthOfQueue) / float64(10)), float64(3)))
        if numberOfNodesToMake > 10 {
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

        env.K8s.ScaleDeployment(deploymentSize, "consumer-deployment")

        delay, err := strconv.Atoi(os.Getenv("CHECK_DELAY"))
        if err == nil {
            fmt.Println("Delaying next check")
            time.Sleep(time.Second * time.Duration(delay))
        } else {
            fmt.Printf("Failed to parse check delay %s", os.Getenv("CHECK_DELAY"))
            time.Sleep(time.Minute * 2)
        }
    }

}
