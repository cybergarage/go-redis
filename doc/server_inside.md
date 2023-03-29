# Inside of go-redis

The go-redis handles RESP (REdis Serialization Protocol) and interprets any commands based on the RESP so that all developers can develop Redis-compatible servers easily. Because the go-redis is a compatible framework based on RESP and Redis command specifications, all developers can connect to a compatible server based on the go-redis using general client drivers for Redis as the following.

<figure>
<img src="img/framework.png" alt="framework" />
</figure>

The go-redis handles all system commands such as PING and SELECT automatically, and so the developers can easily implement their Redis-compatible server only by simply handling user commands such as SET and GET.

## Supported commands

The go-redis defines the [CommandHandler](../redis/handler.go) for the following supported Redis commands as the system and user command handler.

### Connection commands

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: left;">Supported</th>
<th style="text-align: left;">Connection Command</th>
<th style="text-align: left;">Redis Version</th>
<th style="text-align: left;">Note</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>ECHO</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>PING</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>QUIT</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>SELECT</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
</tbody>
</table>

### Server management commands

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: left;">Supported</th>
<th style="text-align: left;">Set Command</th>
<th style="text-align: left;">Redis Version</th>
<th style="text-align: left;">Note</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>CONFIG SET</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>CONFIG GET</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
</tbody>
</table>

### Generic commands

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: left;">Supported</th>
<th style="text-align: left;">Generic Command</th>
<th style="text-align: left;">Redis Version</th>
<th style="text-align: left;">Note</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>COPY</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>DEL</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>DUMP</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>EXISTS</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>EXPIRE</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>EXPIREAT</p></td>
<td style="text-align: left;"><p>1.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>EXPIRETIME</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>KEYS</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>MOVE</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>MIGRATE</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>OBJECT ENCODING</p></td>
<td style="text-align: left;"><p>2.2.3</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>OBJECT FREQ</p></td>
<td style="text-align: left;"><p>4.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>OBJECT HELP</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>OBJECT IDLETIME</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>OBJECT REFCOUNT</p></td>
<td style="text-align: left;"><p>2.2.3</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>PERSIST</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>PEXPIRE</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>PEXPIREAT</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>PEXPIRETIME</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>PTTL</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>RANDOMKEY</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>RENAME</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>RENAMENX</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>RESTORE</p></td>
<td style="text-align: left;"><p>2.8.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SCAN</p></td>
<td style="text-align: left;"><p>2.8.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SORT</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SORT_RO</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>TOUCH</p></td>
<td style="text-align: left;"><p>3.2.1</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>TTL</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>TYPE</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>UNLINK</p></td>
<td style="text-align: left;"><p>4.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>WAIT</p></td>
<td style="text-align: left;"><p>3.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
</tbody>
</table>

### String commands

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: left;">Supported</th>
<th style="text-align: left;">String Command</th>
<th style="text-align: left;">Redis Version</th>
<th style="text-align: left;">Note</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>APPEND</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>DECR</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>DECRBY</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>GET</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>GETDEL</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>GETEX</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>GETRANGE</p></td>
<td style="text-align: left;"><p>2.4.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>GETSET</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>INCR</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>INCRBY</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>INCRBYFLOAT</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>LCS</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>MGET</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>MSET</p></td>
<td style="text-align: left;"><p>1.0.1</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>MSETNX</p></td>
<td style="text-align: left;"><p>1.0.1</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>PSETNX</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>SET</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"><p>Any options are not supported yet</p></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SETEX</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>SETNX</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SERANGE</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>STRLEN</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>SUBSTR</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
</tbody>
</table>

### Hash commands

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: left;">Supported</th>
<th style="text-align: left;">Hash Command</th>
<th style="text-align: left;">Redis Version</th>
<th style="text-align: left;">Note</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HDEL</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HEXISTS</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HGET</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HGETALL</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>HINCRBY</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>HINCRBYFLOAT</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HKEYS</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HLEN</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HMGET</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HMSET</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>HRANDFIELD</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>HSCAN</p></td>
<td style="text-align: left;"><p>2.8.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HSET</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HSETNX</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HSTRLEN</p></td>
<td style="text-align: left;"><p>3.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>HVALS</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
</tbody>
</table>

### List commands

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: left;">Supported</th>
<th style="text-align: left;">List Command</th>
<th style="text-align: left;">Redis Version</th>
<th style="text-align: left;">Note</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BLMOVE</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BLMPOP</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BLPOP</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BRPOP</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BRPOPLPUSH</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>LINDEX</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>LINSERT</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>LLEN</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>LMOVE</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>LMPOP</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>LPOP</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>LPOS</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>LPUSH</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>LPUSHX</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>LRANGE</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>LREM</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>LSET</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>LTRIM</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>RPOP</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>RPOPLPUSH</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>RPUSH</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>RPUSHX</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
</tbody>
</table>

### Set commands

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: left;">Supported</th>
<th style="text-align: left;">Set Command</th>
<th style="text-align: left;">Redis Version</th>
<th style="text-align: left;">Note</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>SADD</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>SCARD</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SDIFF</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SDIFFSTORE</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SINTER</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SINTERCARD</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SINTERSTORE</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>SISMEMBER</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>SMEMBERS</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SMISMEMBER</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SMOVE</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SPOP</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SRANDMEMBER</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>SREM</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SSCAN</p></td>
<td style="text-align: left;"><p>2.8.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SUNION</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SUNIONSTORE</p></td>
<td style="text-align: left;"><p>1.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
</tbody>
</table>

### Sorted set commands

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: left;">Supported</th>
<th style="text-align: left;">Sorted Set Command</th>
<th style="text-align: left;">Redis Version</th>
<th style="text-align: left;">Note</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BZMPOP</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BZPOPMAX</p></td>
<td style="text-align: left;"><p>5.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BZPOPMIN</p></td>
<td style="text-align: left;"><p>5.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>ZADD</p></td>
<td style="text-align: left;"><p>1.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>ZCARD</p></td>
<td style="text-align: left;"><p>1.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZCOUNT</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZDIFF</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZDIFFSTORE</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>ZINCRBY</p></td>
<td style="text-align: left;"><p>1.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZINTER</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZINTERCARD</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZINTERSTORE</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZLEXCOUNT</p></td>
<td style="text-align: left;"><p>2.8.9</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZMPOP</p></td>
<td style="text-align: left;"><p>7.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZMSCORE</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZPOPMAX</p></td>
<td style="text-align: left;"><p>5.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZPOPMIN</p></td>
<td style="text-align: left;"><p>5.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZRANDMEMBER</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>ZRANGE</p></td>
<td style="text-align: left;"><p>1.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZRANGEBYLEX</p></td>
<td style="text-align: left;"><p>2.8.9</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>ZRANGEBYSCORE</p></td>
<td style="text-align: left;"><p>1.0.5</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZRANGESTORE</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZRANK</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>ZREM</p></td>
<td style="text-align: left;"><p>1.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZREMRANGEBYLEX</p></td>
<td style="text-align: left;"><p>2.8.9</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZREMRANGEBYRANK</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZREMRANGEBYSCORE</p></td>
<td style="text-align: left;"><p>1.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>ZREVRANGE</p></td>
<td style="text-align: left;"><p>1.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZREVRANGEBYLEX</p></td>
<td style="text-align: left;"><p>2.8.9</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZREVRANGEBYSCORE</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZREVRANK</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZSCAN</p></td>
<td style="text-align: left;"><p>2.8.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>O</p></td>
<td style="text-align: left;"><p>ZSCORE</p></td>
<td style="text-align: left;"><p>1.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZUNION</p></td>
<td style="text-align: left;"><p>6.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>ZUNIONSTORE</p></td>
<td style="text-align: left;"><p>2.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
</tbody>
</table>

### Bitmap commands

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<thead>
<tr class="header">
<th style="text-align: left;">Supported</th>
<th style="text-align: left;">Bitmap Command</th>
<th style="text-align: left;">Redis Version</th>
<th style="text-align: left;">Note</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BITCOUNT</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BITFIELD</p></td>
<td style="text-align: left;"><p>3.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BITFIELD_RO</p></td>
<td style="text-align: left;"><p>6.0.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BITOP</p></td>
<td style="text-align: left;"><p>2.6.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>BITPOS</p></td>
<td style="text-align: left;"><p>2.8.7</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="even">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>GETBIT</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
<tr class="odd">
<td style="text-align: left;"><p>-</p></td>
<td style="text-align: left;"><p>SETBIT</p></td>
<td style="text-align: left;"><p>2.2.0</p></td>
<td style="text-align: left;"></td>
</tr>
</tbody>
</table>
