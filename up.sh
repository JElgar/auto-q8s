#!/bin/bash

export $(grep -v '^#' .env | xargs -d '\n')
pulumi up -C infra -y
python3 glue_script.py "`pulumi stack -C infra output --json`" --ansible-dir ansible_setup
ansible-playbook -i ansible_setup/inventory.ini ansible_setup/site.yml --user root

export CLUSTER_BASE_URL=$(pulumi stack -C infra output base_url)
export CLUSTER_LOAD_BALANCER_IP=$(pulumi stack -C infra output master-node-0_ip)

cd k8s
./init.sh
