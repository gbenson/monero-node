# Raspberry Pi 4 Monero P2Pool mini node
Ansible playbook to administer a Monero P2Pool mini node on an 8 GB
Raspberry Pi 4 with a 256 GB SD card for the blockchain.  No actual
mining takes place on the Pi, it's far too slow, but cloud-hosting a
P2Pool node needs 8 GB RAM and 70 GB of persistent storage, which will
cost something like $1.60 US a day as of November 2023.  Obviously a
Pi isn't zero-cost either, except I have one doing nothing so it
almost is.

## Setup
Clone the repo:
```sh
git clone https://github.com/gbenson/monero-node.git
cd monero-node/ansible
```
Create a virtual environment:
```sh
python3 -m venv venv
. venv/bin/activate
```
Upgrade pip, and install Ansible:
```sh
pip install --upgrade pip
pip install ansible
```

## Usage
Run the entire playbook:
```sh
ansible-playbook main.yml
```
Skip gathering facts, run only tasks tagged with "services":
```sh
ANSIBLE_GATHERING=explicit ansible-playbook -t services main.yml
```
