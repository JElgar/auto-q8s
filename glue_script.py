import os
import json
import configparser

import infra.constants as constants
import click

import yaml


def write_ansible_config(inventory: dict, ansible_dir: str):
    config = configparser.ConfigParser(allow_no_value=True)
    for group, ips in inventory.items():
        config[group] = {ip: None for ip in ips}

    with open(f"{ansible_dir}/inventory.ini", "w") as configfile:
        config.write(configfile)


def write_ansible_variables(control_plane_endpoint, ansible_dir: str):
    data = {"control_plane_endpoint": control_plane_endpoint}

    with open(f"{ansible_dir}/infra.yml", "w") as file:
        yaml.dump(data, file)


@click.command()
@click.argument("pulumi_state", nargs=1)
@click.option("--ansible-dir", required=True)
def main(pulumi_state, ansible_dir):
    pulumi_state = json.loads(pulumi_state)
    ips = [
        pulumi_state[f"{constants.master_node_name_prefix}_{i}_ip"]
        for i in range(constants.number_of_master_nodes)
    ]
    master_init_ip = ips.pop(0)
    write_ansible_config(
        {
            "master_init_node": [master_init_ip],
            "master_join_nodes": ips,
        },
        ansible_dir,
    )
    write_ansible_variables(pulumi_state["control_plane_endpoint"], ansible_dir)


if __name__ == "__main__":
    main()
