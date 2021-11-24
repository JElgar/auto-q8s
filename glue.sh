#!/bin/bash

pulumi up -C infra -y
python3 create_inventory.py "`pulumi stack -C infra output --json`" --ansible-dir ansible_setup
