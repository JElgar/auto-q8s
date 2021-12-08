package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// "time"
	// "io/ioutil"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/nu7hatch/gouuid"
	// "github.com/sfreiberg/simplessh"
	// "golang.org/x/crypto/ssh"
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

    log.Printf("Join command: ")
    log.Printf(joinCommand)
    options := hcloud.ServerCreateOpts{
        Name: fmt.Sprintf("worker-node-%s", uuid),
        Image: &hcloud.Image{Name: "ubuntu-20.04"},
        ServerType: &hcloud.ServerType{Name: "cx11"},
        SSHKeys: sshKeys,
        UserData: fmt.Sprintf("#cloud-config\nruncmd:\n- touch test-cloudinit.txt\n- curl -s https://raw.githubusercontent.com/JElgar/auto-q8s/main/apps/scaler/init_worker.sh -o init.sh\n- chmod +x init.sh\n- ./init.sh\n- echo 'kubeadm join k8s.james5.joebryan.uk:6443 --token token --discovery-token-ca-cert-hash sha256:something' > join.sh\n- chmod +x join.sh\n- ./join.sh > join_output.txt"), 
    }
    response, _, err := hetzner.Client.Server.Create(context.Background(), options)
    if err != nil {
        log.Printf("Error creating node: %s", err)
        return
    }

    log.Println("Node created")
    log.Println(response)
    // InitNode(response, joinCommand)
  
    // Wait till server is started 
    for {
        action := hetzner.GetAction(response.Action.ID)
        log.Println("Command: ")
        log.Println(action.Command)
        if action.Status == hcloud.ActionStatusRunning {
            log.Println("Running")
        } else if action.Status == hcloud.ActionStatusError {
            log.Println("Error")
            log.Println(action.ErrorMessage)
            log.Println(action.ErrorCode)
            break
        } else if action.Status == hcloud.ActionStatusSuccess {
            log.Println("Success")
            break
        } else {
            log.Println("Unknown status")
        }
        time.Sleep(time.Second * 10)
    }  
}

func (hetzner *Hetzner) GetAction(actionId int) *hcloud.Action {
    action, response, err := hetzner.Client.Action.GetByID(context.Background(), actionId)
    if err != nil {
        log.Println(err)
        log.Println(response)
        log.Println("Failed to get action")
    }
    return action
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
