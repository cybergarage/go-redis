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

## User commands

|Supported|Command|Redis Version|Note                             |
|---------|-------|-------------|---------------------------------|
|O        |GET    |1.0.0        |                                 |
|O        |GETSET |1.0.0        |                                 |
|O        |SET    |1.0.0        |Any options are not supported yet|
|O        |SETNX  |1.0.0        |                                 |
