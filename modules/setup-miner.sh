#!/bin/bash

export PS4='setup-miner.sh: '
set -x

if [ -f /lib/systemd/system/p2pool.service ]; then
  : # called from setup-pool.sh
elif [ -f /etc/amazon-linux-release ]; then
  # Amazon Linux 2023
  dnf install -y docker
  systemctl enable docker
elif grep -q Ubuntu /etc/lsb-release 2>/dev/null; then
  # Ubuntu 22.04
  for snap in lxd core20 snapd; do
    snap remove $snap
  done
  systemctl stop snapd.{service,socket} snapd.{seeded,apparmor}.service
  apt-get autoremove -y --purge snapd
  rm -rf /root/snap

  apt-get autoremove -y --purge modemmanager

  apt-get update
  #apt-get upgrade -y
  apt-get install -y docker.io
fi

secret_dir=/etc/tor-miner
secret_file=config_passphrase
secret_src=$secret_dir/$secret_file
secret_dst=/run/secrets/tor_miner_$secret_file

if [ ! -f $secret_src ]; then
  token=$(curl -X PUT http://169.254.169.254/latest/api/token \
	       -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
  region=$(curl -H "X-aws-ec2-metadata-token: $token" \
		http://169.254.169.254/latest/meta-data/placement/region)

  python3 -m venv /root/venv
  /root/venv/bin/pip install --upgrade pip
  /root/venv/bin/pip install boto3

  mkdir -p $secret_dir
  script='import json,boto3;print(json.loads(boto3.client(service_na'
  script="${script}me='secretsmanager',region_name='$region').get_se"
  script="${script}cret_value(SecretId='tor-miner')['SecretString'])"
  script="${script}['$secret_file'])"
  /root/venv/bin/python -c "$script" > $secret_src
  chmod 600 $secret_src
fi

service=xmrig
cat <<EOF >/lib/systemd/system/$service.service
[Unit]
Description=XMRig Monero miner
Documentation=https://github.com/gbenson/monero-node/
Requires=docker.service
After=docker.service

[Service]
Type=simple
ExecStart=docker run \\
    --privileged \\
    --pull=always \\
    --rm \\
    --name=$service \\
    --mount type=tmpfs,target=/run \\
    --mount type=bind,source=$secret_src,target=$secret_dst,readonly \\
    gbenson/$service
ExecStop=docker stop $service
Restart=always
RestartSec=30s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable $service

if [ -f /lib/systemd/system/p2pool.service ]; then
  : # called from setup-pool.sh
elif [ -f /etc/amazon-linux-release ]; then
  systemctl start $service
else
  reboot
fi
