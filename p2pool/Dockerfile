FROM debian:12 as base

ENV NAME=p2pool
ENV VERSION=v3.7
ENV PACKAGE=$NAME-$VERSION
ENV ARCH=linux-x64
ENV TARBALL=$PACKAGE-$ARCH.tar.gz
ENV DOWNLOAD=https://github.com/SChernykh/$NAME/releases/download
ENV HASH=d5b932053f56f85b119450e64b1a7765c632f16c6226c8762b00e48f50e7e421
ENV UID=802
ENV GID=$UID

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
RUN set -eux \
  \
  && addgroup --system --gid $GID p2pool \
  && adduser --system --uid $UID --gid $GID --disabled-password p2pool

VOLUME /srv
USER p2pool
WORKDIR /srv/p2pool

ENTRYPOINT ["p2pool"]
CMD ["--host",  "monerod", \
     "--mini", \
     "--wallet", "4APvq9TNgrBSYjZ4AvVJpvgh493nzMjmzUNK6EL4aAJiGr1dSGQXx1o5Kk14F4Gj1zDwdfyuBbSwDFKmmFhovTrX976BRNe"]
