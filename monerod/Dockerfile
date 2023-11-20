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
ARG VERSION=v0.18.3.1
ARG DOWNLOAD=https://downloads.getmonero.org/cli

RUN set -eux; \
  \
  case $(uname -m) in \
  x86_64) \
    ARCH=linux-x64 \
    HASH=23af572fdfe3459b9ab97e2e9aa7e3c11021c955d6064b801a27d7e8c21ae09d \
    ;; \
  aarch64) \
    ARCH=linux-armv8 \
    HASH=445032e88dc07e51ac5fff7034752be530d1c4117d8d605100017bcd87c7b21f \
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
