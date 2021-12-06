package services

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
    "strconv"

	"golang.org/x/crypto/ssh"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/nu7hatch/gouuid"
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

func CreateSSHSession(
    host net.IP,
    user string,
    password string,
    port int,
) (*ssh.Session) {
	conf := &ssh.ClientConfig{
		User:            user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    var conn *ssh.Client

    hostWithPort := fmt.Sprintf("%s:22", host.String())
    conn, err := ssh.Dial("tcp", hostWithPort, conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()

    session, err := conn.NewSession()
	if err != nil {
		fmt.Println(err.Error())
    }

    return session
}

func CopyInitScriptToNode(
    host net.IP,
    user string,
    password string,
    port int,
) {
    workingPath, _ := os.Getwd()
    filePath := fmt.Sprintf("%s/scaler/init_worker.sh", workingPath)
    fmt.Println("File path")
    fmt.Println(filePath)


    fmt.Println("Connection host")
    fmt.Println(fmt.Sprintf("%s:%s", host.String(), strconv.Itoa(port)))
    client := scp.NewClient(
        fmt.Sprintf("%s:22", host.String()), 
        &ssh.ClientConfig{
		    User: user,
		    Auth: []ssh.AuthMethod{
			    ssh.Password(password),
		    },
            HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        },
    )
    err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
    }

    fmt.Println("Created client")

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil { 
        fmt.Printf("Could not connect to server: %s", err)
    }

    defer client.Close()
    defer file.Close()

	fmt.Println("Copying file")
    err = client.CopyFile(file, "/test.sh", "0655")
    if err != nil {
		fmt.Println("Error while copying file ", err)
		fmt.Println(err)
    }
	fmt.Println("Copied")
}

func InitNode(response hcloud.ServerCreateResult) {
    time.Sleep(time.Minute * 2)
    CopyInitScriptToNode(
        response.Server.PublicNet.IPv4.IP,
        "root",
        response.RootPassword,
        22,
    )

    session := CreateSSHSession(
        response.Server.PublicNet.IPv4.IP,
        "root",
        response.RootPassword,
        22,
    )

    var b bytes.Buffer  // import "bytes"
    session.Stdout = &b // get output
    // you can also pass what gets input to the stdin, allowing you to pipe
    // content from client to server
    //      session.Stdin = bytes.NewBufferString("My input")

    // Finally, run the command
    err := session.Run("./init.sh")

    fmt.Printf("Stuff happened: %s, %s", b.String(), err)
}

func (hetzner *Hetzner) CreateNode() {
    uuid, _ := uuid.NewV4()
    options := hcloud.ServerCreateOpts{
        Name: fmt.Sprintf("worker-node-%s", uuid),
        Image: &hcloud.Image{Name: "ubuntu-20.04"},
        ServerType: &hcloud.ServerType{Name: "cx11"},
    }
    response, _, err := hetzner.Client.Server.Create(context.Background(), options)
    if err != nil {
        log.Printf("Error creating node: %s", err)
    } else {
        log.Println("Node created")
    }
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
