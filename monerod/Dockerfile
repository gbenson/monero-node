FROM debian:12-slim as base
RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install curl bzip2 \
  && apt-get -y clean \
  && rm -rf /var/lib/apt /var/cache/apt/archives \
  && rm -f /var/cache/ldconfig/aux-cache \
  && for i in /run/*; do [ "$i" = "/run/lock" ] || rm -rf "$i"; done

FROM base as builder

ARG NAME=monero
ARG VERSION=v0.18.3.2
ARG DOWNLOAD=https://downloads.getmonero.org/cli

RUN set -eux; \
  \
  case $(uname -m) in \
  x86_64) \
    ARCH=linux-x64 \
    HASH=9dafd70230a7b3a73101b624f3b5f439cc5b84a19b12c17c24e6aab94b678cbb \
    ;; \
  aarch64) \
    ARCH=linux-armv8 \
    HASH=72f5c90955a736d99c1a645850984535050ebddd42c39a27eec1df82bd972126 \
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

HEALTHCHECK --interval=30s --timeout=5s CMD curl --fail http://localhost:18081/get_info || exit 1

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
