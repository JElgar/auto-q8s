package main

import (
	"apps/services"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
    "strconv"
    "io/ioutil"
    "golang.org/x/crypto/ssh"
)

type Env struct {
    Rmq *services.Rmq
    Hetzner *services.Hetzner
}

func main() {

    // for {
    //     key, err := ioutil.ReadFile("/etc/ssh-key/private-key")
    //     if err != nil {
    //         log.Println(err)
    //         log.Fatalln("Failed to open private key file")
    //     }

    //     signer, err := ssh.ParsePrivateKey(key)
	//     if err != nil {
	//     	log.Fatalf("unable to parse private key: %v", err)
	//     }
    //     fmt.Println("Done signer")

	//     config := &ssh.ClientConfig{
	//     	User: "root",
	//     	Auth: []ssh.AuthMethod{
	//     		// Add in password check here for moar security.
	//     		ssh.PublicKeys(signer),
	//     	},
	//     	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    //     }


    //     client, err := ssh.Dial("tcp", "65.108.147.7:22", config)
    //     if err != nil {
    //         log.Println(err)
	//     	log.Fatal("unable to dial", err)
    //     }
    //     defer client.Close()

	//     session, err := client.NewSession()
	//     if err != nil {
	//     	log.Fatal("unable to create SSH session: ", err)
    //     }
    //     defer session.Close()
    //     fmt.Println("Created session")

    //     session.Run("ls"); 
    //     fmt.Println("Ran init")
    //     panic("All went well")
    // }

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
