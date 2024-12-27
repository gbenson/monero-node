FROM gbenson/tor-node AS base
FROM base AS builder-base
RUN apt-get -y update

FROM builder-base AS xmrig-src

ARG NAME=xmrig
ARG VERSION=6.22.1
ARG TARBALL=v$VERSION.tar.gz
ARG DOWNLOAD=https://github.com/xmrig/xmrig/archive/refs/tags
ARG HASH=189e4bf604f8b569cd15894439a4ad0209b64184c8ccb5a2119a4f07c16da448

RUN apt-get -y install curl
WORKDIR /usr/src
RUN curl -Lo $TARBALL $DOWNLOAD/v$VERSION/$TARBALL
RUN sha256sum $TARBALL
RUN ls -l
RUN echo "$HASH  $TARBALL" | sha256sum -c
RUN tar xvf $TARBALL
RUN mv $NAME-$VERSION $NAME

FROM base AS xmrig-builder
RUN apt-get -y update
RUN apt-get -y install build-essential cmake automake libtool autoconf
RUN apt-get -y install wget
COPY --from=xmrig-src --chown=0:0 /usr/src/xmrig /usr/src/xmrig
WORKDIR /usr/src/xmrig/build
RUN cd ../scripts && ./build_deps.sh
RUN cmake .. \
  -DXMRIG_DEPS=scripts/deps \
  -DWITH_OPENCL=OFF \
  -DWITH_CUDA=OFF \
  -DBUILD_STATIC=ON

RUN make -j$(nproc)

FROM golang:bookworm AS tmn-builder
WORKDIR /usr/src/tor-miner
COPY go.* ./
RUN go mod download
COPY *.go sealed.config ./
COPY cmd ./cmd
RUN go build -trimpath ./cmd/tor-miner

FROM base
RUN set -eux \
  \
  && mkdir /usr/libexec/tor-miner \
  && mv /entrypoint.sh /usr/libexec/tor-miner/start-tor

COPY --from=xmrig-builder --chown=0:0 /usr/src/xmrig/build/xmrig /usr/bin
COPY --from=tmn-builder --chown=0:0 /usr/src/tor-miner/tor-miner /usr/bin

ENTRYPOINT ["tor-miner"]
