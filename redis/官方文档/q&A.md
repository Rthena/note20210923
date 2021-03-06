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

### 3、缓存穿透，缓存击穿，缓存雪崩

### 缓存穿透

**描述：** key对应的数据在数据库不存在，客户端发起请求，请求会先到达缓存，如果缓存不存在，则请求数据库，如果请求大很大，数据源有可能被压垮。

**解决方案：**

- 接口层增加校验，如用户权鉴校验，参数校验， id<=0拦截 。
- 布隆过滤器，将所有的数据哈希到一个足够大的bitmap中，一个不存在的数据会被这个bitmap拦截掉。
- 直接返回带有过期时间的空结果，key-null

### 缓存击穿

**描述：**缓存击穿是指缓存中没有数据但是数据库中有数据（一般是缓存时间到期），这时由于并发用户特别多，引起数据库压力瞬间增大，有可能导致 down 机

**解决方案**

- 设置热点数据永远不过期
- 接口限流与熔断，降级
- 布隆过滤器。bloomfilter就类似于一个hash set，用于快速判某个元素是否存在于集合中，其典型的应用场景就是快速判断一个key是否存在于某容器，不存在就直接返回。布隆过滤器的关键就在于hash算法和容器大小
- 加互斥锁

### 缓存雪崩

**描述：**缓存雪崩是指缓存中数据大批量过期，而查询数据量巨大，引起数据库压力过大导致 down机。和缓存击穿不同的是，缓存击穿指并发查同一条数据，缓存雪崩时是不同 key 都过期。

**解决方案**

- 缓存数据的过期时间设置随机，防止同一时间大量数据过期现象发生。
- 如果缓存数据库是分布式部署，将热点数据均匀分布在不同的缓存数据库中。
- 设置热点数据永远不过期。

