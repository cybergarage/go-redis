# Changelog

## v1.5.0 (2024-05-17)
- New features
  - Supported TLS connection and setting options

## v1.4.5 (2024-05-09)
###  Supported
- Thanks for @Peter-Sh
- SET options and SETEX command

## v1.4.4 (2024-03-20)
- Fix ling warnings

## v1.4.3 (2024-01-26)
- Updated glob package to match more strictly
###  Supported
- SCAN

## v1.4.2 (2023-12-25)
- Added DatabaseID type

## v1.4.1 (2023-12-14)
- Updated Conn interface to embed net.Conn for authenticator

## v1.4.0 (2023-11-20)
- New features
  - Added authorization interface
    - Supported AUTH command
- Changed user command handler interfaces
  - String command handler to remove MGet and MSet interfaces
  - Hash command handler to remove HMGet and HMSet interfaces
- Updated go-tracing package

## v1.3.6 (2023-05-04)
- Updated Conn to embed tracer context
- Updated tracer spans

## v1.3.5 (2023-05-04)
- Updated tracing spans
- Updated Conn interfaces

## v1.3.4 (2023-05-04)
- Updated logger functions to output more detail messages

## v1.3.3 (2023-04-26)
- Add Server::SetPort() and SetAddress()

## v1.3.2 (2023-04-23)
- Add tracing interface

## v1.3.1 (2023-04-02)
- Add connection logs
- Add Dockerfile

## v1.3.0 (2023-03-28)
- Updated executer methods to pass redis.Conn intead of redis.Context
- Add sync.Map interface to redis.Conn to store user data
- Add profiling option to go-redisd

## v1.2.1 (2023-02-24)
- Upgrade to go 1.20
- Updated public helper functions in redistest 

## v1.2.0 (2023-01-13)
- Add a new interface for server management commands
- Updated go-redisd using sync.Map for redis-benchmark
- Tested go-redisd working only with GET/SET commands of redis-benchmark
###  Supported
- CONFIG SET, CONFIG GET

## v1.1.1 (2023-01-02)
- Updated go.mod to go 1.19
- Updated logger package

## v1.1.0 (2022-09-03)
- Support major set, sorted set and list commands
- Enable a test using YCSB (Yahoo! Cloud Serving Benchmark)
###  Supported
- SADD, SCARD, SISMEMBER, SMEMBERS, SREM
- ZADD, ZCARD, ZINCRBY, ZRANGE, ZRANGEBYSCORE, ZREM, ZREVRANGE, ZREVRANGEBYSCORE
- LINDEX, LLEN, LPOP, LPUSH, LPUSHX, LRANGE, RPOP, RPUSH, RPUSHX

## v1.0.0 (2022-08-21)
- Support major generic, string and hash commands
###  Supported
- DEL, EXISTS, EXPIRE, EXPIREAT, KEYS, RENAME, RENAMENX, TTL, TYPE
- APPEND, DECR, DECRBY, GETRANGE, INCR, INCRBY, MGET, MSET, MSETNX, STRLEN, SUBSTR
- HDEL, HEXISTS, HGET, HGETALL, HKEYS, HLEN, HMGET, HMSET, HSET, HSETNX, HSTRLEN, HVALS

## v0.9.0 (2022-07-31)
- Initial release  
###  Supported
- PING, ECHO, SELECT, QUIT
- SET, GET, GETSET, SETNX
