# go-redis

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-redis)
[![test](https://github.com/cybergarage/go-redis/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-redis/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-redis.svg)](https://pkg.go.dev/github.com/cybergarage/go-redis) [![codecov](https://codecov.io/gh/cybergarage/go-redis/graph/badge.svg?token=L7OQDIRHW8)](https://codecov.io/gh/cybergarage/go-redis)

The go-redis is a database framework for implementing a [Redis](https://redis.io)-compatible server using Go easily.

## What is the go-redis?

The go-redis handles RESP (REdis Serialization Protocol) and interprets any commands based on the RESP so that all developers can develop Redis-compatible servers easily. Because the go-redis is a compatible framework based on RESP and Redis command specifications, all developers can connect to a compatible server based on the go-redis using general client drivers for Redis as the following.

![](doc/img/framework.png)

Sinse the go-redis handles all system commands such as PING and SELECT automatically, developers can easily implement their Redis-compatible server only by simply handling user commands such as SET and GET.

## Table of Contents

- [Getting Started](doc/server_impl.md)
- [Inside of go-redis](doc/server_inside.md)

## Examples

- [Examples](doc/examples.md)
  - [go-redisd](examples/go-redisd) [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/go-redisd)](https://hub.docker.com/repository/docker/cybergarage/go-redisd/)
  - [PuzzleDB](https://github.com/cybergarage/puzzledb-go) [![Docker Image Version](https://img.shields.io/docker/v/cybergarage/puzzledb)](https://hub.docker.com/repository/docker/cybergarage/puzzledb/)

# Related Projects

The go-redis is developed in collaboration with the following Cybergarage projects:

-   [go-authenticator](https://github.com/cybergarage/go-authenticator) ![go authenticator](https://img.shields.io/github/v/tag/cybergarage/go-authenticator)
-   [go-logger](https://github.com/cybergarage/go-logger) ![go logger](https://img.shields.io/github/v/tag/cybergarage/go-logger)
-   [go-tracing](https://github.com/cybergarage/go-tracing) ![go tracing](https://img.shields.io/github/v/tag/cybergarage/go-tracing)

## References

- [Redis](https://redis.io)
  - [RESP (REdis Serialization Protocol)](https://github.com/cybergarage/go-redis.git)
- [Rust vs Go vs C: Database and IoT Application Performance Benchmarks – CyberGarage](https://www.cybergarage.org/blog/rust-eval-loc-perf/)