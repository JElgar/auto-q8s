package main

import (
	"apps/services"
	"fmt"
    "strings"
)

func main() {
  hetzner := services.HetznerSetup()

  for _, node := range hetzner.GetNodes() {
    fmt.Println(node)
    if strings.Contains(node.Name, "worker") {
      hetzner.DeleteNode(node.ID)
    }
  }
}
