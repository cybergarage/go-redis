FROM golang:1.20-alpine

USER root

RUN apk add bash

COPY . /go-redis
WORKDIR /go-redis

RUN go build -o /go-redisd github.com/cybergarage/go-redis/examples/go-redisd

ENTRYPOINT ["/go-redisd"]