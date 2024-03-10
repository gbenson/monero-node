FROM debian:12-slim as base
COPY bookworm-backports.list /etc/apt/sources.list.d
RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install --no-install-recommends tor/bookworm-backports \
  && apt-get -y clean \
  && rm -rf /var/lib/apt /var/cache/apt/archives \
  && rm -f /var/cache/ldconfig/aux-cache \
  && for i in /run/*; do [ "$i" = "/run/lock" ] || rm -rf "$i"; done

FROM base as builder
COPY entrypoint.sh /

ARG m_USER=debian-tor
ARG m_UID=803
ARG m_GID=$m_UID

RUN set -eux \
  \
  && groupmod --gid $m_GID $m_USER \
  && usermod --uid $m_UID $m_USER \
  && grep '^ExecStartPre' /lib/systemd/system/tor@default.service \
    | sed 's/^[^=]*=//' >> /entrypoint.sh \
  && grep '^ExecStart=' /lib/systemd/system/tor@default.service \
    | sed 's/^[^=]*=/exec /' >> /entrypoint.sh \
  && chmod +x /entrypoint.sh

VOLUME /var/lib/tor

ENTRYPOINT ["/entrypoint.sh"]
