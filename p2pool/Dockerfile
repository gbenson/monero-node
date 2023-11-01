FROM debian:12-slim as base

ENV NAME=p2pool
ENV VERSION=v3.8
ENV PACKAGE=$NAME-$VERSION
ENV ARCH=linux-x64
ENV TARBALL=$PACKAGE-$ARCH.tar.gz
ENV DOWNLOAD=https://github.com/SChernykh/$NAME/releases/download
ENV HASH=019d62a4bb6d3cf22950e07ad9b49158e6fdadc006a7589edf077b338f0baf98

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
  && curl -Lo $TARBALL $DOWNLOAD/$VERSION/$TARBALL \
  && sha256sum $TARBALL \
  && echo "$HASH  $TARBALL" > SHA256SUMS \
  && sha256sum -c SHA256SUMS \
  && tar xf $TARBALL

FROM base
COPY --from=builder /$PACKAGE-$ARCH/p2pool /usr/bin

ARG m_USER=$NAME
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
     "--wallet", "4APvq9TNgrBSYjZ4AvVJpvgh493nzMjmzUNK6EL4aAJiGr1dSGQXx1o5Kk14F4Gj1zDwdfyuBbSwDFKmmFhovTrX976BRNe"]
