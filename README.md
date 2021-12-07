# auto-q8s

## Setup

Requirements:

- python3.9
- pulumi
- kubectl
- istioctl

Create a .env file with the following values:

```
CLOUDFLARE_API_TOKEN= # Token from cloudflare
HCLOUD_TOKEN= # Token from hetzner 
```

To deploy infra, install k8s and deploy resources run

```
make up
```

This will copy the generated kubeconfig file to your ~/.kube/config file. This will then allow you to conrol the cluster using kubectl. To check everything is working run

```
kubectl get nodes -o wide
```

To destroy run

```
make down
```
