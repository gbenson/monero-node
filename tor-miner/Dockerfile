FROM golang:bookworm as builder

ARG NAME=xmrig
ARG VERSION=6.21.0
ARG ARCH=linux-static-x64
ARG TARBALL=$NAME-$VERSION-$ARCH.tar.gz
ARG DOWNLOAD=https://github.com/xmrig/xmrig/releases/download
ARG HASH=c5dc12dbb9bb51ea8acf93d6349d5bc7fe5ee11b68d6371c1bbb098e21d0f685

WORKDIR /src
RUN curl -Lo $TARBALL $DOWNLOAD/v$VERSION/$TARBALL
RUN echo "$HASH  $TARBALL" | sha256sum -c
RUN tar xvf $TARBALL
RUN mkdir /build
RUN mv $NAME-$VERSION /build/xmrig

WORKDIR /build/xmrig
RUN sha256sum -c SHA256SUMS

WORKDIR /build
COPY go.* ./
RUN go mod download
COPY *.go sealed.config ./
COPY cmd ./cmd
RUN go build -trimpath ./cmd/tor-miner

FROM gbenson/tor-node
RUN set -eux \
  \
  && mkdir /usr/libexec/tor-miner \
  && mv /entrypoint.sh /usr/libexec/tor-miner/start-tor

COPY --from=builder --chown=0:0 /build/xmrig/xmrig /usr/bin
COPY --from=builder --chown=0:0 /build/tor-miner /usr/bin

ENTRYPOINT ["tor-miner"]
