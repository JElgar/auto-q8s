#!/bin/bash

pulumi up -C infra -y
python3 glue_script.py "`pulumi stack -C infra output --json`" --ansible-dir ansible_setup --k8s-dir k8s
ansible-playbook -i ansible_setup/inventory.ini ansible_setup/site.yml --user ubuntu
./k8s/init.sh
