# Monero miner
A toy Monero mining setup I've been using as a test load while I
experiment with deployment and orchestration options on different
compute providers.  Basically just something I can leave running
while I figure out config, pricing and hardening on whatever platform
I'm evaluating.

## Initial setup
### Terraform

This repository is set up to store Terraform's state in a submodule
which it accesses via `git clone terraform@terraform:monero-node`.
To make this work you need to add the hostname or IP address of the
server holding that repo to either `/etc/hosts` or to your
`~/.ssh/config`. I did the latter:

```sh
cat >>~/.ssh/config <<EOF
Host terraform
Hostname (you know it)
ForwardX11 no
ForwardAgent no
Compression no
EOF
```

Once that's done you can recursively clone the repo:
```sh
git clone --recursive https://github.com/gbenson/monero-node.git
cd monero-node
```

### OpenStack client
This part is optional unless you want to use the OpenStack client
directly.  Create a Python virtual environment:
```sh
python3 -m venv venv
. venv/bin/activate
```

Upgrade pip, install OpenStack client:
```sh
pip install --upgrade pip
pip install python-openstackclient
```

## Usage
Source OpenStack configuration and credentials:
```sh
. ~/.config/gbenson/secrets/openstack-openrc.sh
```

Update infrastructure to match definition:
```sh
terraform fmt && terraform plan -out=tfplan
terraform apply tfplan
```
