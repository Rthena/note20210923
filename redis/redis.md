### 1、Redis是什么

Redis 是一个开源的，内存数据结构存储，可用于数据库，缓存，消息发送。Redis提供的数据结构有 strings，hashs， lists，sets， sorted sets，并且还提供了这些数据结构的范围查询，hyperloglogs, geospatial indexes, and stream。Redis还内建了备份，lua脚本，LRU淘汰算法，事务，和不同级别的磁盘持久化。Redis 同时还提供了高可用的 Sentinel 哨兵模式和Redis集群自动分区

### 2、Redis为什么会这么快

Redis 需然是单线程模型，但是Redis 完成是基于内存操作，CPU不是Redis的瓶颈，Redis的最可能瓶颈是内存大小或者网络IO。

Redis为什么这么快；

第一：Redis完成基于内存，纯内存操作，速度非常快。
第二：数据结构简单，对数据操作也简单。
第三：采用单线程，避免不必要的上下稳切换和竞争条件，不存在多线程导致的CPU切换，不用去考虑各种锁的问题，不存在加锁释放锁的操作，没有死锁问题导致的性能消耗
第四：使用IO多路复用模型，非阻塞IO

### 3、Redis和Memcached的区别
1、存储方式上：memcache把数据全部保存在内存中，断电会丢失。redis 有持久化策略，使用过程中会有部分数据落盘。
2、数据支持类型：memcache仅支持简单的 key-value的数据类型。redis支持5种数据类型。
3、底层模型不同：
4、value的大小：redis可以达到1GB，而memcache只有1MB

### 4、[淘汰策略](https://redis.io/topics/lru-cache)
当达到设置的最大内存使用值 maxmemory时采用的淘汰策略
**noevition**：不淘汰策略，当内存达到了限定值时，客户端试图执行那些会导致内存使用更多的命令时，会之间报错（大多数是写操作命令，del命令和少数例外）
**allkeys-lru**：从所以keys 中淘汰最近最少使用的（LRU）keys
**volatile-lru**：从已设置了过期时间的keys 优先淘汰最近最少使用的数据。
**allkeys-random**：从所有的keys 中随机淘汰
**volatile-random**：从已设置过期时间的keys 中随机淘汰
**volatile-ttl**：从已设置过期时间的keys 中优先对剩余时间短的(ttl)的数据淘汰

### 5、[持久化机制](https://redis.io/topics/persistence)
1、RDB：
