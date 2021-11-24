#!/bin/bash

pulumi up -y
python3 create_inventory.py "`pulumi stack output --json`"
