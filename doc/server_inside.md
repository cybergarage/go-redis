# Inside of go-redis

The go-redis handles RESP (REdis Serialization Protocol) and interprets any commands based on the RESP so that all developers can develop Redis compatible servers easily. Because the go-redis is a compatible framework based on RESP and Redis command specifications, all developers can connect to a compatible server based on the go-redis using general client drivers for Redis as the following.

![](img/framework.png)

The go-redis handles all system commands such as PING and SELECT automatically, and so the developers can easily implement their Redis compatible server only by simply handling user commands such as SET and GET.


## Supported commands

## System commands

|Supported|Command|Redis Version|Note                             |
|---------|-------|-------------|---------------------------------|
|O        |ECHO   |1.0.0        |                                 |
|O        |PING   |1.0.0        |                                 |
|O        |QUIT   |1.0.0        |                                 |
|O        |SELECT |1.0.0        |                                 |

## Generic commands

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

## String commands

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

## Hash commands

|Supported|Hash Command      |Redis Version|Note                             |
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
