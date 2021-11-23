## Pulumi 

Creates 3 nodes

## Ansible config
We want to generate an inventory

```
[master_nodes]
10.1.6.10
10.1.6.11
10.1.6.12
```

## Install
Code in the string in the dev thing
Source:
https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/

## Init
https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/

## Glue
pulumi up
python convert_to_inventory `pulumi output --json` -> Create ansible thing from output
ansible go -i theoutputfromabove.txt


# Autoscaling
https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/clusterapi/README.md

