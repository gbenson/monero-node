FROM debian:12-slim as base

ENV NAME=monero
ENV VERSION=v0.18.3.1
ENV TARBALL=$NAME-linux-x64-$VERSION.tar.bz2
ENV DOWNLOAD=https://downloads.getmonero.org/cli
ENV HASH=23af572fdfe3459b9ab97e2e9aa7e3c11021c955d6064b801a27d7e8c21ae09d

RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install curl bzip2 \
  && apt-get -y clean \
  && rm -rf /var/lib/apt

FROM base as builder
RUN set -eux \
  \
  && curl -Lo $TARBALL $DOWNLOAD/$TARBALL \
  && sha256sum $TARBALL \
  && echo "$HASH  $TARBALL" > SHA256SUMS \
  && sha256sum -c SHA256SUMS \
  && tar xf $TARBALL \
  && rm $TARBALL \
  && mv /$NAME-* /$NAME

FROM base
COPY --from=builder /$NAME/$NAME* /usr/bin/

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
