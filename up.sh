#!/bin/bash

pulumi up -C infra -y
python3 create_inventory.py "`pulumi stack -C infra output --json`" --ansible-dir ansible_setup
ansible-playbook -i ansible_setup/inventory.ini ansible_setup/site.yml --user ubuntu