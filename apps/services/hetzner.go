package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/nu7hatch/gouuid"
    "github.com/sfreiberg/simplessh"
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

func InitNode(response hcloud.ServerCreateResult) {
    time.Sleep(time.Second * 20)
    var client *simplessh.Client
    var err error

    host := fmt.Sprintf("%s:22", response.Server.PublicNet.IPv4.IP.String())
    homedir, err := os.UserHomeDir()
    if client, err = simplessh.ConnectWithKeyFile(host, "root", fmt.Sprintf("%s/.ssh/id_rsa", homedir); err != nil {
        fmt.Println(err)
        return
    }
    defer client.Close()

    // Now run the init script from github 
    _, err = client.Exec("bash <(curl -s https://raw.githubusercontent.com/JElgar/auto-q8s/main/apps/scaler/init_worker.sh)"); 
    if err != nil {
        log.Println(err)
    }
}

func (hetzner *Hetzner) getSshKeyId() (*hcloud.SSHKey, error) {
    options := hcloud.SSHKeyListOpts{
        Name: "jelgar@JamesLaptop",
    }

    response, err := hetzner.Client.SSHKey.AllWithOpts(context.Background(), options)
    if err != nil || len(response) == 0 {
        fmt.Println("Could not get ssh key")
        return nil, err
    }
    return response[0], nil
}

func (hetzner *Hetzner) CreateNode() {
    uuid, _ := uuid.NewV4()
    sshKey, err := hetzner.getSshKeyId()
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
    InitNode(response)
}

func (hetzner *Hetzner) DeleteNode(name string) {
    server := &hcloud.Server{
        Name: name,
    }
    _, err := hetzner.Client.Server.Delete(context.Background(), server)
    if err != nil {
        log.Printf("Error deleting node: %s", err)
    } else {
        log.Println("Node deleted")
    }
}
