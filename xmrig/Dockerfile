FROM debian:12 as base

ENV NAME=xmrig
ENV VERSION=6.20.0
ENV PACKAGE=$NAME-$VERSION
ENV ARCH=linux-static-x64
ENV TARBALL=$PACKAGE-$ARCH.tar.gz
ENV DOWNLOAD=https://github.com/xmrig/$NAME/releases/download
ENV HASH=ff6e67d725ee64b4607dc6490a706dc9234c708cff814477de52d3beb781c6a1

RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install curl \
  && apt-get -y clean \
  && rm -rf /var/lib/apt

FROM base as builder
RUN set -eux \
  \
  && curl -Lo $TARBALL $DOWNLOAD/v$VERSION/$TARBALL \
  && sha256sum $TARBALL \
  && echo "$HASH  $TARBALL" > SHA256SUMS \
  && sha256sum -c SHA256SUMS \
  && tar xf $TARBALL \
  && cd $NAME-$VERSION \
  && sha256sum -c SHA256SUMS

FROM base
RUN mkdir -m0700 /etc/xmrig
COPY --from=builder /$NAME-$VERSION/xmrig /usr/bin
COPY --from=builder --chmod=600 /$NAME-$VERSION/config.json /etc/xmrig

VOLUME /srv
WORKDIR /srv/xmrig

ENTRYPOINT ["xmrig"]
CMD ["-c", "/etc/xmrig/config.json", "-o",  "p2pool:3333"]
