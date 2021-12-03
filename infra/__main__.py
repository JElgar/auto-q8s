"""An AWS Python Pulumi program"""

from pathlib import Path

import pulumi
import pulumi_aws as aws
from pulumi_aws.ec2.key_pair import KeyPair
from pulumi_aws.ec2.security_group import SecurityGroup
import pulumi_cloudflare as cloudflare
import pulumi_hcloud as hcloud

import constants


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


def create_aws_node(
    name: str, ami: str, ssh_key: KeyPair, security_group: SecurityGroup
) -> aws.ec2.Instance:
    instance = aws.ec2.Instance(
        name,
        ami=ami,
        instance_type="t2.medium",
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


def create_hetzner_node(
    name: str,
    ssh_key_name: str,
) -> hcloud.Server:
    instance = hcloud.Server(
        name, image="ubuntu-20.04", ssh_keys=[ssh_key_name], server_type="cx21"
    )

    # Export server details
    pulumi.export(f"{name}_ip", instance.ipv4_address)

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


def create_dns_record(resource_name: str, name: str, target) -> cloudflare.Record:
    return cloudflare.Record(
        resource_name,
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
    pulumi.export("control_plane_endpoint", f"{k8s_subdomain}.{constants.zone_domain}")
    for i in range(number_of_master_nodes):
        node = create_aws_node(
            name=f"{constants.master_node_name_prefix}_{i}",
            ami=ubuntu_ami,
            ssh_key=key,
            security_group=security_group,
        )
        if i == 0:
            create_dns_record(f"{k8s_subdomain}_{i}", k8s_subdomain, node.public_ip)


def create_master_nodes_hetzner(number_of_master_nodes: int) -> None:
    key = "jelgar@JamesLaptop"

    k8s_subdomain = f"k8s.{pulumi.get_stack()}"
    pulumi.export("control_plane_endpoint", f"{k8s_subdomain}.{constants.zone_domain}")
    for i in range(number_of_master_nodes):
        node = create_hetzner_node(
            name=f"{constants.master_node_name_prefix}-{i}",
            ssh_key_name=key,
        )
        if i == 0:
            create_dns_record(f"{k8s_subdomain}_{i}", k8s_subdomain, node.ipv4_address)


pulumi.export("base_url", f"{pulumi.get_stack()}.{constants.zone_domain}")
# create_master_nodes(constants.number_of_master_nodes)
create_master_nodes_hetzner(constants.number_of_master_nodes)
