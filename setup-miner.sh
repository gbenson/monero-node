#!/bin/bash

export PS4='setup-miner.sh: '
set -ex

name=xmrig
version=6.20.0
package=$name-$version
arch=linux-static-x64
tarball=$package-$arch.tar.gz
download=https://github.com/xmrig/$name/releases/download
hash=ff6e67d725ee64b4607dc6490a706dc9234c708cff814477de52d3beb781c6a1

if [ -d /opt/xmrig ]; then
  SKIP_SETUP=yes
else
  for snap in lxd core20 snapd; do
    snap remove $snap
  done
  systemctl stop snapd.{service,socket} snapd.{seeded,apparmor}.service
  apt-get autoremove -y --purge snapd
  rm -rf /root/snap

  apt-get autoremove -y --purge modemmanager

  #apt-get update
  #apt-get upgrade -y
  #apt-get dist-upgrade -y
fi

mkdir -p /opt/xmrig/bin && cd /opt/xmrig
curl -Lo $tarball $download/v$version/$tarball
echo $hash $tarball > $tarball.SHA256SUM
sha256sum -c $tarball.SHA256SUM
tar xf $tarball
mv $name-$version $version
cd $version
sha256sum -c SHA256SUMS
ln -sf $(pwd)/xmrig /opt/xmrig/bin

if [ ! -f /lib/systemd/system/xmrig.service ]; then
  cat <<EOF >/lib/systemd/system/xmrig.service
[Unit]
Description=XMRig Monero miner
Documentation=https://github.com/gbenson/monero-node
After=network.target

[Service]
Type=simple
ExecStart=/opt/xmrig/bin/xmrig -o p2pool.gbenson.net:3333
Restart=always
RestartSec=30s

[Install]
WantedBy=multi-user.target
EOF

  systemctl daemon-reload
  systemctl enable xmrig
fi

[ -n "$SKIP_SETUP" ] || reboot
