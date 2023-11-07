# Monero miner
A toy Monero mining setup I've been using as a test load while I
experiment with deployment and orchestration options on different
compute providers.  Basically just something I can leave running
while I figure out config, pricing and hardening on whatever platform
I'm evaluating.

## Setup
Clone the repo:
```sh
git clone https://github.com/gbenson/monero-node.git
cd monero-node
```
Create a virtual environment:
```sh
python3 -m venv venv
. venv/bin/activate
```
Upgrade pip, install OpenStack client:
```sh
pip install --upgrade pip
pip install python-openstackclient
```
Source OpenStack config and credentials:
```sh
. ~/.config/gbenson/secrets/openstack-openrc.sh
```

## Usage
```sh
terraform fmt && terraform plan -out=tfplan
terraform apply tfplan
```
