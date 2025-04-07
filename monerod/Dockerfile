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

RUN --mount=type=tmpfs,target=/var/cache \
    --mount=type=tmpfs,target=/var/lib/apt/lists \
    --mount=type=tmpfs,target=/var/log \
  set -eux \
  \
  && apt-get -y update \
  && apt-get -y install --no-install-recommends \
	bzip2 \
	ca-certificates \
	curl \
  && rm -f /var/lib/dpkg/*-old

FROM builder-base AS builder

ARG NAME=monero
ARG VERSION=v0.18.4.0
ARG DOWNLOAD=https://downloads.getmonero.org/cli

RUN set -eux; \
  \
  case $(uname -m) in \
  x86_64) \
    ARCH=linux-x64 \
    HASH=16cb74c899922887827845a41d37c7f3121462792a540843f2fcabcc1603993f \
    ;; \
  aarch64) \
    ARCH=linux-armv8 \
    HASH=f252b6a24e801535bf36fbaaa7b2d6ae44b1efc5d427803d483e3c3a17d6f2cd \
    ;; \
  *) exit 1 ;; \
  esac; \
  TARBALL=$NAME-$ARCH-$VERSION.tar.bz2; \
  echo "$HASH  $TARBALL" > SHA256SUMS \
  && curl -Lo $TARBALL $DOWNLOAD/$TARBALL \
  && sha256sum -c SHA256SUMS \
  && tar xf $TARBALL \
  && mv /$NAME-*-$VERSION /build

FROM base
COPY --from=builder /build/monero* /usr/bin/

ARG m_USER=monerod
ARG m_UID=801
ARG m_GID=$m_UID

RUN set -eux \
  \
  && addgroup --system --gid $m_GID $m_USER \
  && adduser --system --uid $m_UID --gid $m_GID --disabled-password $m_USER \
  && mkdir -m 0700 /srv/monerod \
  && chown $m_UID:$m_GID /srv/monerod

VOLUME /srv/monerod

USER $m_USER:$m_USER

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
