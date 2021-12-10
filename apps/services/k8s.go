package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
)

type K8sEnv struct {
	Clientset *kubernetes.Clientset
}

func int32Ptr(i int32) *int32 { return &i }

func InitK8s() *K8sEnv {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &K8sEnv {
		Clientset: clientset,
	}
}

func (env K8sEnv) GetNodes() []apiv1.Node {
	nodes, err := env.Clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))
	return nodes.Items
}


func (env K8sEnv) GetNodeByName(name string) *apiv1.Node {
	nodes := env.GetNodes()
	for _, node := range nodes {
		if node.Name == name {
			return &node 
		}
	}
	return nil
}


func (env K8sEnv) GetWorkerNodes() []apiv1.Node {
	nodes := env.GetNodes() 

	var workerNodes[]apiv1.Node

    for _, node := range nodes {
		if strings.Contains(node.Name, "worker") {
			workerNodes = append(workerNodes, node)
		}
    }
		
	return workerNodes
}

func (env K8sEnv) GetNodeCount() int {
	return len(env.GetNodes())
}

func (env K8sEnv) getNodeIp(node apiv1.Node) string {
	nodeAddress := node.Status.Addresses
	for _, address := range nodeAddress {
		if address.Type == apiv1.NodeInternalIP {
			return address.String()
		}
	}
	return ""
}

func (env K8sEnv) DeleteNodes (numberOfNodes int, hetzner *Hetzner) {
	workerNodes := env.GetWorkerNodes()[:numberOfNodes]
	for _, node := range workerNodes {
		env.DeleteNode(node, hetzner)
	}
}

func (env K8sEnv) DeleteNode(node apiv1.Node, hetzner *Hetzner) {
	log.Printf("Deleting node: %s", node.Name)
	env.Clientset.CoreV1().Nodes().Delete(context.Background(), node.Name, metav1.DeleteOptions{})
	hetznerNode := hetzner.GetNodeByName(node.Name)
	if hetznerNode == nil {
		log.Println("Could not find hetzner node!")
	}
	hetzner.DeleteNode(hetznerNode.ID)
	log.Println("Node deleted")
}

func (env K8sEnv) ScaleDeployment(num_replicas int, deployment string) {
	deploymentsClient := env.Clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(context.TODO(), deployment, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(int32(num_replicas))
		_, updateErr := deploymentsClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}
