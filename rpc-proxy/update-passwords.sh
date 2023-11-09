#!/bin/sh

srcdir=/var/lib/docker/volumes/monerod/_data
dstdir=$(cd $(dirname "$0") && pwd)

for line in $(grep 'rpc-login[[:space:]]*=' $srcdir/bitmonero.conf); do
  pair=$(echo $line | sed 's/^[^=]*=[[:space:]]*//')
  username=$(echo $pair | awk -F: '{ print $1 }')
  password=$(echo $pair | awk -F: '{ print $2 }')
  prefix="$username:monero-rpc:"
  digest=$(echo -n "$prefix$password" | md5sum | sed 's/[[:space:]].*$//')
  echo "$prefix$digest"
done > $dstdir/.htpasswd.digest
