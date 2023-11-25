#!/bin/bash

export PS4='setup-miner.sh: '
set -ex

if [ -f /etc/amazon-linux-release ]; then
  # Amazon Linux 2023
  dnf install docker
  systemctl enable docker
else
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

service=xmrig
secret_name=config_passphrase
secret_src=/etc/tor-miner/$secret_name
secret_dst=/run/secrets/tor_miner_${secret_name}

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

reboot
