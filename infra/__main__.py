"""An AWS Python Pulumi program"""

from pathlib import Path

import pulumi
import pulumi_aws as aws
from pulumi_aws.ec2.key_pair import KeyPair
from pulumi_aws.ec2.security_group import SecurityGroup
import pulumi_cloudflare as cloudflare

import constants

# https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
install_kubeadm = """
#!/bin/bash

sudo apt-get update -y
sudo apt-get install -y ec2-instance-connect

cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sudo sysctl --system

sudo apt-get install -y containerd

sudo apt-get update -y
sudo apt-get install -y apt-transport-https ca-certificates curl

sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

sudo apt-get update -y
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
"""


def get_ubuntu_ami() -> str:
    return aws.ec2.get_ami(
        most_recent=True,
        filters=[
            aws.ec2.GetAmiFilterArgs(
                name="name",
                values=["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"],
            ),
            aws.ec2.GetAmiFilterArgs(
                name="virtualization-type",
                values=["hvm"],
            ),
        ],
        owners=["099720109477"],
    ).id


def create_key():
    homedir = str(Path.home())
    sshkey_path = f"{homedir}/.ssh/id_rsa.pub"
    sshkey_file = open(sshkey_path, "r")
    return aws.ec2.KeyPair(
        f"{pulumi.get_stack()}_key", public_key=sshkey_file.read().strip("\n")
    )


def create_node(
    name: str, ami: str, ssh_key: KeyPair, security_group: SecurityGroup
) -> aws.ec2.Instance:
    instance = aws.ec2.Instance(
        name,
        ami=ami,
        instance_type="t2.micro",
        user_data=install_kubeadm,
        tags={
            "stack": pulumi.get_stack(),
            "Name": name,
        },
        key_name=ssh_key.key_name,
        vpc_security_group_ids=[security_group.id],
    )

    # Export server details
    pulumi.export(f"{name}_arn", instance.arn)
    pulumi.export(f"{name}_ip", instance.public_ip)

    return instance


def create_security_group():
    group = aws.ec2.SecurityGroup(
        f"{pulumi.get_stack()}_master_node_security_group",
        description="Enable all tcp access",
        ingress=[
            aws.ec2.SecurityGroupIngressArgs(
                protocol="tcp",
                from_port=0,
                to_port=65535,
                cidr_blocks=["0.0.0.0/0"],
            )
        ],
        egress=[
            aws.ec2.SecurityGroupEgressArgs(
                protocol="tcp",
                from_port=0,
                to_port=65535,
                cidr_blocks=["0.0.0.0/0"],
            )
        ],
    )
    return group


def create_dns_record(name: str, target) -> cloudflare.Record:
    return cloudflare.Record(
        name,
        zone_id=constants.zone_id,
        name=name,
        value=target,
        type="A",
        ttl=3600,
    )


def create_master_nodes(number_of_master_nodes: int) -> None:
    ubuntu_ami = get_ubuntu_ami()
    key = create_key()
    security_group = create_security_group()

    k8s_subdomain = f"k8s.{pulumi.get_stack()}"
    pulumi.export("control_plane_endpoint", k8s_subdomain)
    for i in range(number_of_master_nodes):
        node = create_node(
            name=f"{constants.master_node_name_prefix}_{i}",
            ami=ubuntu_ami,
            ssh_key=key,
            security_group=security_group,
        )
        create_dns_record(k8s_subdomain, node.public_ip)


create_master_nodes(constants.number_of_master_nodes)
