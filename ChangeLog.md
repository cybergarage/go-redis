# Changelog

## v1.3.3 (2023-04-26)
- Add Server::SetPort() and SetAddress()

## v1.3.2 (2023-04-23)
- Add tracing interface

## v1.3.1 (2023-04-02)
- Add connection logs
- Add Dockerfile

## v1.3.0 (2023-03-28)
- Update executer methods to pass redis.Conn intead of redis.Context
- Add sync.Map interface to redis.Conn to store user data
- Add profiling option to go-redis-server

## v1.2.1 (2023-02-24)
- Upgrade to go 1.20
- Update public helper functions in redistest 

## v1.2.0 (2023-01-13)
- Add a new interface for server management commands
- Update go-redis-server using sync.Map for redis-benchmark
- Check go-redis-server working only with GET/SET commands of redis-benchmark
###  Supported
- CONFIG SET, CONFIG GET

## v1.1.1 (2023-01-02)
- Update go.mod to go 1.19
- Update logger package

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
