# Inside of go-redis

The go-redis handles RESP (REdis Serialization Protocol) and interprets any commands based on the RESP so that all developers can develop Redis-compatible servers easily. Because the go-redis is a compatible framework based on RESP and Redis command specifications, all developers can connect to a compatible server based on the go-redis using general client drivers for Redis as the following.

![](img/framework.png)

The go-redis handles all system commands such as PING and SELECT automatically, and so the developers can easily implement their Redis-compatible server only by simply handling user commands such as SET and GET.


## Supported commands

The go-redis defines the [CommandHandler](../redis/handler.go) for the following supported Redis commands as the system and user command handler.

### System commands

|Supported|Command|Redis Version|Note                             |
|---------|-------|-------------|---------------------------------|
|O        |ECHO   |1.0.0        |                                 |
|O        |PING   |1.0.0        |                                 |
|O        |QUIT   |1.0.0        |                                 |
|O        |SELECT |1.0.0        |                                 |

### Generic commands

|Supported|Command      |Redis Version|Note                             |
|---------|------------------|-------------|---------------------------------|
|-        |COPY              |6.2.0        |                                 |
|O        |DEL               |1.0.0        |                                 |
|-        |DUMP              |2.6.0        |                                 |
|O        |EXISTS            |1.0.0        |                                 |
|O        |EXPIRE            |1.0.0        |                                 |
|O        |EXPIREAT          |1.2.0        |                                 |
|-        |EXPIRETIME        |7.0.0        |                                 |
|O        |KEYS              |1.0.0        |                                 |
|-        |MOVE              |1.0.0        |                                 |
|-        |MIGRATE           |2.6.0        |                                 |
|-        |RANDOMKEY         |1.0.0        |                                 |
|O        |RENAME            |1.0.0        |                                 |
|O        |RENAMENX          |1.0.0        |                                 |
|-        |RESTORE           |2.8.0        |                                 |
|-        |SCAN              |2.8.0        |                                 |
|-        |SORT              |1.0.0        |                                 |
|-        |SORT_RO           |7.0.0        |                                 |
|-        |TOUCH             |3.2.1        |                                 |
|O        |TTL               |1.0.0        |                                 |
|O        |TYPE              |1.0.0        |                                 |
|-        |UNLINK            |4.0.0        |                                 |
|-        |WAIT              |3.0.0        |                                 |

### String commands

|Supported|Command      |Redis Version|Note                             |
|---------|------------------|-------------|---------------------------------|
|O        |APPEND            |2.0.0        |                                 |
|O        |DECR              |1.0.0        |                                 |
|O        |DECRBY            |1.0.0        |                                 |
|O        |GET               |1.0.0        |                                 |
|-        |GETDEL            |6.2.0        |                                 |
|-        |GETEX             |6.2.0        |                                 |
|O        |GETRANGE          |2.4.0        |                                 |
|O        |GETSET            |1.0.0        |                                 |
|O        |INCR              |1.0.0        |                                 |
|O        |INCRBY            |1.0.0        |                                 |
|-        |INCRBYFLOAT       |2.6.0        |                                 |
|-        |LCS               |7.0.0        |                                 |
|O        |MGET              |1.0.0        |                                 |
|O        |MSET              |1.0.1        |                                 |
|O        |MSETNX            |1.0.1        |                                 |
|-        |PSETNX            |2.6.0        |                                 |
|O        |SET               |1.0.0        |Any options are not supported yet|
|-        |SETEX             |2.0.0        |                                 |
|O        |SETNX             |2.0.0        |                                 |
|-        |SERANGE           |2.2.0        |                                 |
|O        |STRLEN            |2.2.0        |                                 |
|O        |SUBSTR            |1.0.0        |                                 |

### Hash commands

|Supported|Command      |Redis Version|Note                             |
|---------|------------------|-------------|---------------------------------|
|O        |HDEL              |2.0.0        |                                 |
|O        |HEXISTS           |2.0.0        |                                 |
|O        |HGET              |2.0.0        |                                 |
|O        |HGETALL           |2.0.0        |                                 |
|-        |HINCRBY           |2.0.0        |                                 |
|-        |HINCRBYFLOAT      |2.6.0        |                                 |
|O        |HKEYS             |2.0.0        |                                 |
|O        |HLEN              |2.0.0        |                                 |
|O        |HMGET             |2.0.0        |                                 |
|O        |HMSET             |2.0.0        |                                 |
|-        |HRANDFIELD        |6.2.0        |                                 |
|-        |HSCAN             |2.8.0        |                                 |
|O        |HSET              |2.0.0        |                                 |
|O        |HSETNX            |2.0.0        |                                 |
|O        |HSTRLEN           |3.2.0        |                                 |
|O        |HVALS             |2.0.0        |                                 |

### List commands

|Supported|\Command|Redis Version|Note|
|---------|------------|-------------|----|
|-        |BLMOVE      |7.0.0        |    |
|-        |BLMPOP      |7.0.0        |    |
|-        |BLPOP       |2.0.0        |    |
|-        |BRPOP       |2.0.0        |    |
|-        |BRPOPLPUSH  |2.2.0        |    |
|O        |LINDEX      |1.0.0        |    |
|-        |LINSERT     |2.2.0        |    |
|O        |LLEN        |1.0.0        |    |
|-        |LMOVE       |6.2.0        |    |
|-        |LMPOP       |7.0.0        |    |
|O        |LPOP        |1.0.0        |    |
|-        |LPOS        |6.2.0        |    |
|O        |LPUSH       |1.0.0        |    |
|O        |LPUSHX      |2.2.0        |    |
|O        |LRANGE      |1.0.0        |    |
|-        |LREM        |1.0.0        |    |
|-        |LSET        |1.0.0        |    |
|-        |LTRIM       |1.0.0        |    |
|O        |RPOP        |1.0.0        |    |
|-        |RPOPLPUSH   |6.2.0        |    |
|O        |RPUSH       |1.0.0        |    |
|O        |RPUSHX      |2.2.0        |    |

### Set commands

|Supported|Command|Redis Version|Note|
|---------|-----------|-------------|----|
|O        |SADD       |1.0.0        |    |
|O        |SCARD      |1.0.0        |    |
|-        |SDIFF      |1.0.0        |    |
|-        |SDIFFSTORE |1.0.0        |    |
|-        |SINTER     |1.0.0        |    |
|-        |SINTERCARD |7.0.0        |    |
|-        |SINTERSTORE|1.0.0        |    |
|O        |SISMEMBER  |1.0.0        |    |
|O        |SMEMBERS   |1.0.0        |    |
|-        |SMISMEMBER |6.2.0        |    |
|-        |SMOVE      |1.0.0        |    |
|-        |SPOP       |1.0.0        |    |
|-        |SRANDMEMBER|1.0.0        |    |
|O        |SREM       |1.0.0        |    |
|-        |SSCAN      |2.8.0        |    |
|-        |SUNION     |1.0.0        |    |
|-        |SUNIONSTORE|1.0.0        |    |

### Sorted set commands

|Supported|Command|Redis Version|Note|
|---------|------------------|-------------|----|
|-        |BZMPOP            |7.0.0        |    |
|-        |BZPOPMAX          |5.0.0        |    |
|-        |BZPOPMIN          |5.0.0        |    |
|O        |ZADD              |1.2.0        |    |
|O        |ZCARD             |1.2.0        |    |
|-        |ZCOUNT            |2.0.0        |    |
|-        |ZDIFF             |6.2.0        |    |
|-        |ZDIFFSTORE        |6.2.0        |    |
|O        |ZINCRBY           |1.2.0        |    |
|-        |ZINTER            |6.2.0        |    |
|-        |ZINTERCARD        |7.0.0        |    |
|-        |ZINTERSTORE       |2.0.0        |    |
|-        |ZLEXCOUNT         |2.8.9        |    |
|-        |ZMPOP             |7.0.0        |    |
|-        |ZMSCORE           |6.2.0        |    |
|-        |ZPOPMAX           |5.0.0        |    |
|-        |ZPOPMIN           |5.0.0        |    |
|-        |ZRANDMEMBER       |6.2.0        |    |
|O        |ZRANGE            |1.2.0        |    |
|-        |ZRANGEBYLEX       |2.8.9        |    |
|O        |ZRANGEBYSCORE     |1.0.5        |    |
|-        |ZRANGESTORE       |6.2.0        |    |
|-        |ZRANK             |2.0.0        |    |
|O        |ZREM              |1.2.0        |    |
|-        |ZREMRANGEBYLEX    |2.8.9        |    |
|-        |ZREMRANGEBYRANK   |2.0.0        |    |
|-        |ZREMRANGEBYSCORE  |1.2.0        |    |
|O        |ZREVRANGE         |1.2.0        |    |
|-        |ZREVRANGEBYLEX    |2.8.9        |    |
|-        |ZREVRANGEBYSCORE  |2.2.0        |    |
|-        |ZREVRANK          |2.0.0        |    |
|-        |ZSCAN             |2.8.0        |    |
|O        |ZSCORE            |1.2.0        |    |
|-        |ZUNION            |6.2.0        |    |
|-        |ZUNIONSTORE       |2.0.0        |    |
