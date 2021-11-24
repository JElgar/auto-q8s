import json
import configparser

import constants
import click

import yaml


def write_ansible_config(inventory: dict):
    config = configparser.ConfigParser(allow_no_value=True)
    for group, ips in inventory.items():
        config[group] = {ip: None for ip in ips}

    with open("inventory.ini", "w") as configfile:
        config.write(configfile)


def write_ansible_variables(control_plane_subdomain):
    data = [
        {"control_plane_enpoint": f"{control_plane_subdomain}.{constants.zone_domain}"}
    ]

    with open("infra.yml", "w") as file:
        yaml.dump(data, file)


@click.command()
@click.argument("pulumi_state", nargs=1)
def main(pulumi_state):
    data = json.loads(pulumi_state)
    ips = [
        data[f"{constants.master_node_name_prefix}_{i}_ip"]
        for i in range(constants.number_of_master_nodes)
    ]
    write_ansible_config({"master_nodes": ips})
    write_ansible_variables(data["control_plane_endpoint"])


if __name__ == "__main__":
    main()
