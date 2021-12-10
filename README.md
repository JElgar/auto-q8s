# auto-q8s

## Setup

Requirements:

- python3.9
- pulumi
- kubectl
- istioctl
- ansible

For development:
- docker
- docker-compose

Create a .env file with envrionment varaible matching those in .env.example

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

## Services

Defore deploying any services you must set your environment varaibles, run the following at the root of the project.

```
export $(grep -v '^#' .env | xargs -d '\n')
export CLUSTER_BASE_URL=$(pulumi stack -C infra output base_url)
export CLUSTER_LOAD_BALANCER_IP=$(pulumi stack -C infra output master-node-0_ip)
export JOIN_COMMAND=$(cat ansible_setup/join_command)
export DYNAMO_TABLE=$(pulumi stack -C infra output table_name)
```

### Rabbitmq

To deploy a queue to the cluster go into k8s/rabbitmq and run the init.sh script

### Scaler 

To deploy the scaler to the cluster go into apps/scaler/infra and run the init.sh script

### Producer

To deploy the producer to the cluster go into apps/producer/infra and run the init.sh script

### Consumer 

To deploy the cosumer to the cluster go into apps/consumer/infra and run the init.sh script

### Metrics 

To deploy the metrics service to the cluster go into apps/metrics/infra and run the init.sh script
