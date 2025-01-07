FROM debian:12-slim AS base
RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y clean \
  && rm -rf /var/lib/apt /var/cache/apt/archives \
  && rm -f /var/cache/ldconfig/aux-cache \
  && for i in /run/*; do [ "$i" = "/run/lock" ] || rm -rf "$i"; done

FROM base AS builder

ARG NAME=p2pool
ARG VERSION=v4.3
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
  HASH=$(grep -A2 $TARBALL $SHA256SUMS | tail -1 | sed 's/.*\([0-9A-F]\{64\}\).*/\1/'); \
  echo "$HASH  $TARBALL" > SHA256SUMS \
  && curl -Lo $TARBALL $DOWNLOAD/$VERSION/$TARBALL \
  && sha256sum -c SHA256SUMS \
  && tar xf $TARBALL \
  && mv /$PACKAGE-$ARCH /build

FROM base
COPY --from=builder /build/p2pool /usr/bin

ARG m_USER=p2pool
ARG m_UID=802
ARG m_GID=$m_UID

RUN set -eux \
  \
  && addgroup --system --gid $m_GID $m_USER \
  && adduser --system --uid $m_UID --gid $m_GID --disabled-password $m_USER \
  && mkdir -m 0700 /srv/p2pool \
  && chown $m_UID:$m_GID /srv/p2pool

VOLUME /srv/p2pool

USER $m_USER:$m_USER
WORKDIR /srv/p2pool

ENTRYPOINT ["p2pool"]
CMD ["--host",  "monerod", \
     "--mini", \
     "--no-upnp", \
     "--wallet", "4APvq9TNgrBSYjZ4AvVJpvgh493nzMjmzUNK6EL4aAJiGr1dSGQXx1o5Kk14F4Gj1zDwdfyuBbSwDFKmmFhovTrX976BRNe"]
