FROM debian:12-slim

COPY entrypoint.sh /
RUN set -eux \
  \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install tor \
  && apt-get -y clean \
  && rm -rf /var/lib/apt \
  && sed -i \
    -e '0,/^#HiddenServicePort/{s/^#\(HiddenService\)/\1/}' \
    -e '/^HiddenServicePort/{s/127\.0\.0\.1/p2pool/; s/80/3333/g}' \
    /etc/tor/torrc \
  && grep '^ExecStartPre' /lib/systemd/system/tor@default.service \
    | sed 's/^[^=]*=//' >> /entrypoint.sh \
  && grep '^ExecStart=' /lib/systemd/system/tor@default.service \
    | sed 's/^[^=]*=/exec /' >> /entrypoint.sh \
  && chmod +x /entrypoint.sh

VOLUME /var/lib/tor

CMD /entrypoint.sh
