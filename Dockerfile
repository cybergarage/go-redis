FROM alpine:latest
RUN apk update && apk add git go

USER root

COPY . /go-redis
WORKDIR /go-redis

RUN go build -o /go-redisd github.com/cybergarage/go-redis/examples/go-redisd

ENTRYPOINT ["/go-redisd"]
