FROM golang:1.19-alpine3.16 as builder

WORKDIR /reloader
COPY --chown=root:root ./reloader .
RUN go build

FROM golang:1.19-alpine3.16 as runner

ENV GOPATH=/go
ENV APPDIR=/app
ENV GOSRC=${GOPATH}/src
ENV PKGSITE_PORT=3000
ENV PROXY_PORT=80

EXPOSE ${PROXY_PORT}

RUN go install github.com/yosssi/goat@v0.0.0-20190705092005-4e07e5bfb19f
RUN go install golang.org/x/pkgsite/cmd/pkgsite@v0.0.0-20220825124633-4a62ba3611bc

COPY --from=builder --chown=root:root /reloader/reloader /usr/local/bin/reloader

WORKDIR ${APPDIR}
COPY --chown=root:root ./app/* .
RUN ln start.sh /usr/local/bin/start
RUN ln startservice.sh /usr/local/bin/startservice
RUN ln stopservice.sh /usr/local/bin/stopservice
WORKDIR ${GOSRC}

CMD [ "start" ]
