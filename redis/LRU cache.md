### 4、[淘汰策略](https://redis.io/topics/lru-cache)
当达到设置的最大内存使用值 maxmemory时采用的淘汰策略
**noevition**：不淘汰策略，当内存达到了限定值时，客户端试图执行那些会导致内存使用更多的命令时，会之间报错（大多数是写操作命令，del命令和少数例外）
**allkeys-lru**：从所以keys 中淘汰最近最少使用的（LRU）keys
**volatile-lru**：从已设置了过期时间的keys 优先淘汰最近最少使用的数据。
**allkeys-random**：从所有的keys 中随机淘汰
**volatile-random**：从已设置过期时间的keys 中随机淘汰
**volatile-ttl**：从已设置过期时间的keys 中优先对剩余时间短的(ttl)的数据淘汰
