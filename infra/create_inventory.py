import json
import configparser

import constants
import click

def write_ansible_config(inventory: dict):
    config = configparser.ConfigParser(allow_no_value=True)
    for group, ips in inventory.items():
        config[group] = {ip: None for ip in ips}

    with open('inventory.ini', 'w') as configfile:
        config.write(configfile)


@click.command()
@click.argument("pulumi_state", nargs=1)
def main(pulumi_state):
    data = json.loads(pulumi_state)
    ips = [
        data[f"{constants.master_node_name_prefix}_{i}_ip"]
        for i in range(constants.number_of_master_nodes)
    ]
    write_ansible_config({"master_nodes": ips})


if __name__ == "__main__":
    main()
