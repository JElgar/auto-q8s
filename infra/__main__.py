"""An AWS Python Pulumi program"""

import pulumi
from pulumi_aws import s3
import pulumi_aws as aws

number_of_master_nodes = 3

ubuntu = aws.ec2.get_ami(
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
)

for i in range(number_of_master_nodes):
    instance = aws.ec2.Instance(
        f"web_{i}",
        ami=ubuntu.id,
        instance_type="t2.micro",
        tags={
            "stack": pulumi.get_stack(),
            "Name": f"web_{i}",
        },
    )

    # Export server details
    pulumi.export(f"master_{i}_arn", instance.arn)
    pulumi.export(f"master_{i}_ip", instance.public_ip)
