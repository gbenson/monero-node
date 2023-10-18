# https://github.com/sethforprivacy/simple-monerod-docker
# https://hub.docker.com/r/sethsimmons/simple-monerod/tags
FROM sethsimmons/simple-monerod:v0.18.3.1

ARG m_USER=monerod
ARG m_UID=801
ARG m_GID=$m_UID
ARG m_HOME=/nonexistent
ARG m_SHELL=/usr/sbin/nologin

USER root
WORKDIR /
RUN set -eux \
  \
  && sed -i "s/fixuid -q//" /entrypoint.sh \
  && rm -f /usr/local/bin/fixuid \
  && chown 0:0 /usr/local/bin/monerod \
  && mv /usr/local/bin/monerod /usr/bin \
  && deluser --remove-home monero \
  && addgroup -S -g $m_GID $m_USER \
  && adduser -h $m_HOME -s $m_SHELL -G $m_USER -u $m_UID -D -g "" $m_USER \
  && mkdir -m 0700 /srv/monerod \
  && chown $m_UID:$m_GID /srv/monerod

VOLUME /srv/monerod

USER $m_USER:$m_USER

CMD ["--data-dir", "/srv/monerod", \
     "--max-log-file-size", "10485760", \
     "--max-log-files", "4", \
     "--rpc-bind-ip", "0.0.0.0", "--confirm-external-bind", \
     "--zmq-pub", "tcp://0.0.0.0:18083", \
     "--no-igd", \
     "--out-peers", "16", "--in-peers", "8", \
     "--add-priority-node=p2pmd.xmrvsbeast.com:18080", \
     "--add-priority-node=nodes.hashvault.pro:18080", \
     "--disable-dns-checkpoints"]
