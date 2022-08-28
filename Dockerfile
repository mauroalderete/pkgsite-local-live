FROM golang:1.19-alpine3.16 as builder

ENV WORKDIR=/app

RUN apk add inotify-tools>3
COPY --chown=root:root ./watcher/watcher.sh /usr/local/bin/watcher
WORKDIR ${WORKDIR}
COPY --chown=root:root ./app/* .
CMD [ "./start.sh" ]
