FROM debian:12-slim AS base
ENV DEBIAN_FRONTEND=noninteractive

RUN --mount=type=tmpfs,target=/var/cache \
    --mount=type=tmpfs,target=/var/lib/apt/lists \
    --mount=type=tmpfs,target=/var/log \
  set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade --no-install-recommends \
  && rm -f /var/lib/dpkg/*-old

FROM base AS builder-base

RUN apt-get -y update
RUN apt-get -y install \
	bzip2 \
	ca-certificates \
	curl \
	gpg

FROM builder-base AS builder

ARG NAME=monero
ARG VERSION=v0.18.4.0
ARG DOWNLOAD=https://downloads.getmonero.org/cli
ARG KEYRING=binaryfate.asc

WORKDIR /usr/src/$NAME-$VERSION

COPY $KEYRING .
RUN gpg --import $KEYRING

RUN curl -L https://www.getmonero.org/downloads/hashes.txt -o hashes.txt.asc
RUN gpg --output hashes.txt --decrypt hashes.txt.asc

RUN set -eux; \
  case $(uname -m) in \
  x86_64) \
    arch=linux-x64 \
    ;; \
  aarch64) \
    arch=linux-armv8 \
    ;; \
  *) exit 1 ;; \
  esac; \
  \
  grep "[[:space:]]$NAME-$arch-v.*\\.tar\\.bz2\$" hashes.txt > SHA256SUMS \
  && cat -n SHA256SUMS \
  && [ "$(wc -l SHA256SUMS)" = '1 SHA256SUMS' ] \
  && grep "[[:space:]]$NAME-$arch-$VERSION\\.tar\\.bz2\$" hashes.txt \
  || (echo -e >&2 "\033[30;1;48;5;220m version bump? \033[0m"; exit 1)

RUN set -eux; \
  tarball=$(awk '{ print $2 }' SHA256SUMS); \
  \
  curl -L $DOWNLOAD/$tarball -o $tarball \
  && sha256sum -c SHA256SUMS \
  && tar xf $tarball \
  && mv $NAME-*-$VERSION /build

FROM base
COPY --from=builder /build/monero* /usr/bin/

RUN set -eux \
  \
  && addgroup --system --gid 801 monerod \
  && adduser --system --uid 801 --gid 801 --disabled-password monerod \
  && install -d -omonerod -gmonerod -m700 /srv/monerod

VOLUME /srv/monerod

USER monerod:monerod

ENTRYPOINT ["monerod"]
CMD ["--non-interactive", \
     "--data-dir=/srv/monerod", \
     "--log-level=0", \
     "--rpc-bind-ip=0.0.0.0", "--confirm-external-bind", \
     "--zmq-pub=tcp://0.0.0.0:18083", \
     "--no-igd", \
     "--out-peers=8", "--in-peers=16", \
     "--add-priority-node=p2pmd.xmrvsbeast.com:18080", \
     "--add-priority-node=nodes.hashvault.pro:18080", \
     "--disable-dns-checkpoints", \
     "--enable-dns-blocklist"]
