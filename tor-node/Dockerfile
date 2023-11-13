FROM debian:12-slim as base
RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install tor \
  && apt-get -y clean \
  && rm -rf /var/lib/apt

FROM base as builder
COPY entrypoint.sh /
RUN set -eux \
  \
  && grep '^ExecStartPre' /lib/systemd/system/tor@default.service \
    | sed 's/^[^=]*=//' >> /entrypoint.sh \
  && grep '^ExecStart=' /lib/systemd/system/tor@default.service \
    | sed 's/^[^=]*=/exec /' >> /entrypoint.sh \
  && chmod +x /entrypoint.sh

VOLUME /var/lib/tor

ENTRYPOINT ["/entrypoint.sh"]