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

FROM base AS builder

ARG NAME=p2pool
ARG VERSION=v4.5
ARG PACKAGE=$NAME-$VERSION
ARG DOWNLOAD=https://github.com/SChernykh/$NAME/releases/download
ARG SHA256SUMS=sha256sums.txt.asc

RUN apt-get -y update
RUN apt-get -y install curl gpg
RUN curl -Lo $SHA256SUMS $DOWNLOAD/$VERSION/$SHA256SUMS
COPY SChernykh.asc /
RUN gpg --import SChernykh.asc
RUN gpg --verify $SHA256SUMS
RUN set -eux; \
  \
  ARCH=linux-$(uname -m | sed s/86_//); \
  TARBALL=$PACKAGE-$ARCH.tar.gz; \
  HASH=$(grep -A2 $TARBALL $SHA256SUMS | tail -1 | sed 's/.*\([0-9a-fA-F]\{64\}\).*/\1/'); \
  echo "$HASH  $TARBALL" > SHA256SUMS \
  && curl -Lo $TARBALL $DOWNLOAD/$VERSION/$TARBALL \
  && sha256sum -c SHA256SUMS \
  && tar xf $TARBALL \
  && mv /$PACKAGE-$ARCH /build

FROM base
COPY --from=builder /build/p2pool /usr/bin

RUN set -eux \
  \
  && addgroup --system --gid 802 p2pool \
  && adduser --system --uid 802 --gid 802 --disabled-password p2pool \
  && install -d -op2pool -gp2pool -m700 /srv/p2pool

VOLUME /srv/p2pool

USER p2pool:p2pool
WORKDIR /srv/p2pool

ENTRYPOINT ["p2pool"]
CMD ["--host",  "monerod", \
     "--mini", \
     "--no-upnp", \
     "--wallet", "4APvq9TNgrBSYjZ4AvVJpvgh493nzMjmzUNK6EL4aAJiGr1dSGQXx1o5Kk14F4Gj1zDwdfyuBbSwDFKmmFhovTrX976BRNe"]
