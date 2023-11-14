FROM golang:bookworm as builder

ARG NAME=xmrig
ARG VERSION=6.20.0
ARG ARCH=linux-static-x64
ARG TARBALL=$NAME-$VERSION-$ARCH.tar.gz
ARG DOWNLOAD=https://github.com/xmrig/xmrig/releases/download
ARG HASH=ff6e67d725ee64b4607dc6490a706dc9234c708cff814477de52d3beb781c6a1

WORKDIR /src
RUN curl -Lo $TARBALL $DOWNLOAD/v$VERSION/$TARBALL
RUN echo "$HASH  $TARBALL" | sha256sum -c
RUN tar xvf $TARBALL
RUN mkdir /build
RUN mv $NAME-$VERSION /build/xmrig

WORKDIR /build/xmrig
RUN sha256sum -c SHA256SUMS
RUN sed \
  -e 's/"[^"]*:3333"/"@@POOL@@"/' \
  -e 's/YOUR_WALLET_ADDRESS/@@WALLET@@/' \
  config.json >config.json.in

COPY config.sh .
RUN sed \
  -e 's:^\([^=]*\)=\(.*\)$:s/@@\1@@/\2/g:' \
  config.sh > config.sed
RUN sed \
  -f config.sed \
  config.json.in >config.json

WORKDIR /build
COPY go.* *.go ./
COPY cmd ./cmd
RUN go build -trimpath -o tor-miner ./cmd/main.go

FROM gbenson/tor-node
RUN set -eux \
  \
  && mkdir /usr/libexec/tor-miner \
  && mv /entrypoint.sh /usr/libexec/tor-miner/start-tor \
  && mkdir -m0700 /etc/xmrig

COPY --from=builder --chown=0:0 /build/xmrig/xmrig /usr/bin
COPY --from=builder --chown=0:0 /build/xmrig/config.json /etc/xmrig
COPY --from=builder --chown=0:0 /build/tor-miner /usr/bin

ENTRYPOINT ["tor-miner"]
