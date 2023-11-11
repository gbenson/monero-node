FROM debian:12-slim as base
RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install curl \
  && apt-get -y clean \
  && rm -rf /var/lib/apt

FROM base as builder

ARG NAME=xmrig
ARG VERSION=6.20.0
ARG PACKAGE=$NAME-$VERSION
ARG ARCH=linux-static-x64
ARG TARBALL=$PACKAGE-$ARCH.tar.gz
ARG DOWNLOAD=https://github.com/xmrig/$NAME/releases/download
ARG HASH=ff6e67d725ee64b4607dc6490a706dc9234c708cff814477de52d3beb781c6a1

RUN set -eux \
  \
  && echo "$HASH  $TARBALL" > SHA256SUMS \
  && curl -Lo $TARBALL $DOWNLOAD/v$VERSION/$TARBALL \
  && sha256sum -c SHA256SUMS \
  && tar xf $TARBALL \
  && mv /$NAME-$VERSION /build \
  && cd /build \
  && sha256sum -c SHA256SUMS

FROM base
RUN mkdir -m0700 /etc/xmrig
COPY --from=builder /build/xmrig /usr/bin
COPY --from=builder --chmod=600 /build/config.json /etc/xmrig

VOLUME /srv
WORKDIR /srv/xmrig

ENTRYPOINT ["xmrig"]
CMD ["-c", "/etc/xmrig/config.json", "-o",  "p2pool:3333"]
