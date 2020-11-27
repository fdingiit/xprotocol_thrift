## tbase-go-client v1.0.0

### Features
1. support operation asynchronized
2. tbase log's directory is ~/logs/tbase-go-client/tbase-go-client.log

## tbase-go-client v0.0.1

### Features

1. Supports Print like tbase commands;
2. Supports tbase cluster and shard auto routing;
3. Supports error like MOVED/READONLY etc. auto handling;
4. Supports logging.

v0.0.1 only supports single key commands and the supported commands are listed as belows:

 - SET/GET
 - SETTSEX/GETTSEX
 - SETEX
 - TTL
 - DEL  

## tbase-go-client v1.1
### Features
1. Supports allmost all print like single key tbase commands;
2. Supports tbase cluster and shard auto routing;
3. Supports error like MOVED/READONLY etc. auto handling;
4. Supports stable and user ignorant net error processing; 
5. Supports logging;
 
### Missing Features
1. multiple commands;
2. hotkey listening;
3. perf log.