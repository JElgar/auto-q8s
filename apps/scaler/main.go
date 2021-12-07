package main

import (
	"apps/services"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
    "strconv"
	"os/exec"
)

type Env struct {
    Rmq *services.Rmq
    Hetzner *services.Hetzner
}

func main() {

    log.Println("Running test ssh")
    cmd := fmt.Sprintf("ssh -i /etc/ssh-key/private-key %s 'ls -al'", "65.108.147.7")
    out, err := exec.Command("bash", "-c", cmd).Output()
    if err != nil {
        log.Println("Something went wrong")
        log.Println(err)
    }
    fmt.Println(out)
    log.Println("Ran test ssh")

    fmt.Println("Getting hetzner stuff")
    fmt.Println(os.Getenv("HCLOUD_TOKEN"))

    hetzner := services.HetznerSetup();
    sshKey, err := hetzner.GetSshKeyId()

    if err != nil {
        log.Panicln("Cannot get key")
    }
    fmt.Println("Got the id!")
    fmt.Println(sshKey)

    rmq := services.RabbitmqSetup()
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
        numberOfNodesToMake = 1 
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
