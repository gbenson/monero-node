FROM gbenson/tor-node as base
FROM base as builder-base
RUN apt-get -y update

FROM builder-base as xmrig-src

ARG NAME=xmrig
ARG VERSION=6.21.1
ARG TARBALL=v$VERSION.tar.gz
ARG DOWNLOAD=https://github.com/xmrig/xmrig/archive/refs/tags
ARG HASH=c45baea5a210143b647349b5234a2192164d3473a39d2b1cab7fb35b1a2a8ba9

RUN apt-get -y install curl
WORKDIR /usr/src
RUN curl -Lo $TARBALL $DOWNLOAD/v$VERSION/$TARBALL
RUN sha256sum $TARBALL
RUN ls -l
RUN echo "$HASH  $TARBALL" | sha256sum -c
RUN tar xvf $TARBALL
RUN mv $NAME-$VERSION $NAME

FROM base as xmrig-builder
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

FROM golang:bookworm as tmn-builder
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
