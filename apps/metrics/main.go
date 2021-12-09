package main 

import (
	"apps/services"
	"fmt"
	"math"
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

    startTime := time.Now()
        
    for {
        time.Sleep(1)
        lengthOfQueue := env.Rmq.QueueLength()
        currentNumberOfJoinedNodes := env.K8s.GetNodeCount()
        currentNumberOfCloudNodes := len(env.Hetzner.GetNodes())
        deploymentSize := int(math.Max(math.Ceil(float64(lengthOfQueue) / float64(10)), float64(3)))

        fmt.Printf(
            "%f,%d,%d,%d,%d",
            time.Since(startTime).Seconds(),
            lengthOfQueue,
            currentNumberOfJoinedNodes,
            currentNumberOfCloudNodes,
            deploymentSize,
        )
    }
}

