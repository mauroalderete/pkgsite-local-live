FROM golang:1.19-alpine3.16 as builder

ENV GOPATH=/go
ENV APPDIR=/app
ENV GOSRC=${GOPATH}/src

EXPOSE 8080

RUN go install github.com/yosssi/goat@v0.0.0-20190705092005-4e07e5bfb19f
RUN go install golang.org/x/pkgsite/cmd/pkgsite@v0.0.0-20220825124633-4a62ba3611bc

WORKDIR ${APPDIR}
COPY --chown=root:root ./app/* .
RUN ln startwatcher.sh /usr/local/bin/startwatcher
RUN ln startservice.sh /usr/local/bin/startservice
RUN ln stopservice.sh /usr/local/bin/stopservice
WORKDIR ${GOSRC}

CMD [ "startwatcher" ]
