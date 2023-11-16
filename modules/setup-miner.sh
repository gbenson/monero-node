#!/bin/bash

export PS4='setup-miner.sh: '
set -ex

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

cat <<EOF >/lib/systemd/system/xmrig.service
[Unit]
Description=XMRig Monero miner
Documentation=https://github.com/gbenson/monero-node
Requires=docker.service
After=docker.service

[Service]
Type=simple
ExecStart=docker run --privileged --pull=always --rm --name=xmrig gbenson/xmrig
ExecStop=docker stop xmrig
Restart=always
RestartSec=30s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable xmrig

reboot
