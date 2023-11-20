FROM debian:12-slim as base
RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install curl \
  && apt-get -y clean \
  && rm -rf /var/lib/apt /var/cache/apt/archives \
  && rm -f /var/cache/ldconfig/aux-cache \
  && for i in /run/*; do [ "$i" = "/run/lock" ] || rm -rf "$i"; done

FROM base as builder

ARG NAME=p2pool
ARG VERSION=v3.8
ARG PACKAGE=$NAME-$VERSION
ARG DOWNLOAD=https://github.com/SChernykh/$NAME/releases/download

RUN set -eux; \
  \
  case $(uname -m) in \
  x86_64) \
    ARCH=linux-x64 \
    HASH=019d62a4bb6d3cf22950e07ad9b49158e6fdadc006a7589edf077b338f0baf98 \
    ;; \
  aarch64) \
    ARCH=linux-aarch64 \
    HASH=3b233ae9e08fc2d5cd6698a1706c5869c9efa36b6db5e2752b766e29bc599f07 \
    ;; \
  *) exit 1 ;; \
  esac; \
  TARBALL=$PACKAGE-$ARCH.tar.gz; \
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
