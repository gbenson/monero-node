#!/bin/bash
set -e

if [ $# -gt 1 ]; then
  echo 1>&2 "docker run gbenson/tor-node [[HOST:]PORT]"
  exit 1
fi

unset host port
if [ $# -eq 1 ]; then
  port=$(echo "$1" | sed 's/.*://')
  if [ "$port" != "$1" ]; then
    host=$(echo "$1" | sed "s/:$port\$//")
  fi
fi

if [ -n "$port" ]; then
  f=/etc/tor/torrc
  grep ^HiddenServicePort $f && exit 1  # crowbar
  sed -e '0,/^#HiddenServicePort/{s/^#\(HiddenService\)/\1/}' \
      -e '/^HiddenServicePort/{s/80/@@PORT@@/g}' \
      -i $f
  if [ -n "$host" ]; then
    sed -i "/^HiddenServicePort/{s/127\.0\.0\.1/@@HOST@@/}" $f
  fi
  sed -e "/^HiddenServicePort/{s/@@HOST@@/$host/; s/@@PORT@@/$port/g}" \
      -i $f
  unset f
fi

unset host port

chown debian-tor:debian-tor /var/lib/tor
chmod 02700 /var/lib/tor
