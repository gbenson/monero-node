FROM debian:12-slim as base

ENV NAME=p2pool
ENV VERSION=v3.7
ENV PACKAGE=$NAME-$VERSION
ENV ARCH=linux-x64
ENV TARBALL=$PACKAGE-$ARCH.tar.gz
ENV DOWNLOAD=https://github.com/SChernykh/$NAME/releases/download
ENV HASH=d5b932053f56f85b119450e64b1a7765c632f16c6226c8762b00e48f50e7e421

RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install curl netcat-openbsd \
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

ARG hcr='{"id":1,"method":"login","params":{"login":"healthcheck"}}'
ARG hcx='"status":"OK"'
HEALTHCHECK --interval=30s --timeout=5s CMD echo $hcr | nc -W1 127.0.0.1 3333 | grep -q $hcx

ENTRYPOINT ["p2pool"]
CMD ["--host",  "monerod", \
     "--mini", \
     "--wallet", "4APvq9TNgrBSYjZ4AvVJpvgh493nzMjmzUNK6EL4aAJiGr1dSGQXx1o5Kk14F4Gj1zDwdfyuBbSwDFKmmFhovTrX976BRNe"]
