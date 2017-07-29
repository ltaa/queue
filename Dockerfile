FROM golang:1.8-alpine3.6

ENV PROJNAME=queue
RUN mkdir -p /app/src/$PROJNAME
ENV GOPATH=/app
ENV GOBIN=$GOPATH/bin

COPY . /app/src/$PROJNAME
WORKDIR /app

RUN apk add --no-cache libc-dev gcc make git
RUN go-wrapper download $PROJNAME/main
RUN go build -o $GOBIN/server $PROJNAME/main
RUN git clone git://github.com/bmc/daemonize.git && cd daemonize && ./configure && make


FROM alpine:latest
RUN mkdir -p /app/src/$PROJNAME && mkdir -p /usr/local/sbin
RUN apk --no-cache add ca-certificates rsyslog
WORKDIR /app/
COPY --from=0 /app/bin/server .
COPY --from=0 /app/daemonize/daemonize /usr/local/sbin

CMD rsyslogd && daemonize /app/server && sh
