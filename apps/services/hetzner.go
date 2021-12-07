package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/nu7hatch/gouuid"
    // "github.com/sfreiberg/simplessh"
    "golang.org/x/crypto/ssh"
)

type Hetzner struct {
    Client *hcloud.Client
}

func HetznerSetup() *Hetzner {
    fmt.Printf("Hetzner token is %s", os.Getenv("HCLOUD_TOKEN"))
    client := hcloud.NewClient(
        hcloud.WithToken(
            os.Getenv("HCLOUD_TOKEN"),
        ),
    )
    return &Hetzner{Client: client} 
}

func InitNode(response hcloud.ServerCreateResult, joinCommand string) {
    time.Sleep(time.Second * 20)
    fmt.Println("The private key is: ")
    fmt.Println(os.Getenv("SSH_PRIVATE_KEY"))
    key := []byte(os.Getenv("SSH_PRIVATE_KEY"))

    signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
    fmt.Println("Done signer")

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    host := response.Server.PublicNet.IPv4.IP.String()
    port := "22"
    client, err := ssh.Dial("tcp", host+":"+port, config)
    if err != nil {
        log.Println(err)
		log.Fatal("unable to dial", err)
    }

    defer client.Close()
    fmt.Println("Created clinet")

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("unable to create SSH session: ", err)
    }
    defer session.Close()
    fmt.Println("Created session")

    session.Run("bash <(curl -s https://raw.githubusercontent.com/JElgar/auto-q8s/main/apps/scaler/init_worker.sh)"); 
    fmt.Println("Ran init")
   
    session.Run(joinCommand); 
    fmt.Println("Ran join")
}

func (hetzner *Hetzner) GetSshKeyId() (*hcloud.SSHKey, error) {
    options := hcloud.SSHKeyListOpts{
        Name: "jelgar@JamesLaptop",
    }

    response, err := hetzner.Client.SSHKey.AllWithOpts(context.Background(), options)
    if err != nil || len(response) == 0 {
        fmt.Println("Could not get ssh key")
        fmt.Println(response)
        fmt.Println(err)
        return nil, err
    }
    return response[0], nil
}

func (hetzner *Hetzner) CreateNode(joinCommand string) {
    uuid, _ := uuid.NewV4()
    sshKey, err := hetzner.GetSshKeyId()
    if err != nil {
        return
    }
    sshKeys := make([]*hcloud.SSHKey, 1)
    sshKeys[0] = sshKey 

    options := hcloud.ServerCreateOpts{
        Name: fmt.Sprintf("worker-node-%s", uuid),
        Image: &hcloud.Image{Name: "ubuntu-20.04"},
        ServerType: &hcloud.ServerType{Name: "cx11"},
        SSHKeys: sshKeys,
    }
    response, _, err := hetzner.Client.Server.Create(context.Background(), options)
    if err != nil {
        log.Printf("Error creating node: %s", err)
        return
    }

    log.Println("Node created")
    InitNode(response, joinCommand)
}

func (hetzner *Hetzner) GetNodes() ([]*hcloud.Server) {
    nodes, err := hetzner.Client.Server.All(context.Background())
    if err != nil {
        fmt.Println("Failed to get nodes")
        return []*hcloud.Server{}
    }
    return nodes
}

func (hetzner *Hetzner) DeleteNode(id int) {
    server := &hcloud.Server{
        ID: id,
    }
    _, err := hetzner.Client.Server.Delete(context.Background(), server)
    if err != nil {
        log.Printf("Error deleting node: %s", err)
    } else {
        log.Println("Node deleted")
    }
}
