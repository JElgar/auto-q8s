#!bin/bash

pulumi up -y
python create_inventory.py "`pulumi stack output --json`"
