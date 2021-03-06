### 使用Redis 用作 LRU缓存

当 Redis 用作缓存时，通常在添加新数据时让它自动 淘汰 旧数据。这种行为在开发者社区是众所周知的，因为它是流行的 memcached 系统的默认行为。

事实上 LRU也是唯一支持的 淘汰方法。本页涵盖了 Redis maxmemory 指令更一般主题，该指令用于限制内存固定使用数量，并且还深入介绍 Redis 使用的 LRU 算法，这实际上是精确 LRU 的近似值。

### 最大内存配置指令

修改配置文件 redis.conf，或者在运行时使用 CONFIG SET 命令，如下

```
maxmemory 100mb
```

设置 maxmemory = 0时，即是没有内存使用限制。这是 64 位系统的默认行为，而 32 位系统使用 3GB 的隐式内存限制

当内存使用数量达到指定的值时，Redis 可以在不同的行为中选择 ，我们称之为 策略 **policies**。

**淘汰策略（Eviction policies）**

- **noeviction**: 当内存已经到达限制时，客户端任然执行可能会导致更多内存使用的命令时候，直接返回错误（大多数写命令，除来 DEL 和少一下例外）。
- **allkeys-lru**: 从所有的 keys 优先淘汰最近最少使用的数据。
- **volatile-lru**: 从设置了过期时间的 keys 中优先淘汰最近最少使用的数据。
- **allkeys-random**: 从所有的 keys 中随机淘汰数据。
- **volatile-random**: 从设置了过期时间的 keys 中随机淘汰数据。
- **volatile-ttl**: 从设置了过期时间的 keys 中，优先淘汰 剩余时间最短的数据。

如果没有与先决条件匹配的淘汰key **volatile-lru**, **volatile-random** and **volatile-ttl** 的策略行为类似与 **noeviction**

选择正确的淘汰策略重要的是取决于你的应用访问模式，您也可以在运行时重新设置策略，通过使用 Redis 的 **INFO** 输出，监控有多少缓存未命中和命中的数量来调整您的设置

一般来说，作为经验法则

- 使用 **allkeys-lru** 策略，当您期望请求的受欢迎程度呈幂律分布（二八定律），也就是，您期望的元素子集的访问频率远高于其他元素。如果您不确定，这将是一个不错的选择。
- 使用 **allkeys-random** 策略，如果您的访问是循环的，其中所有的 keys 都被连续扫描，或者当您期望分布是均匀的（所有元素可能有相同的访问概率）。
- 使用 **volatile-ttl** 策略 如果您能通过创建缓存对象时，设置不同的 TTL 值，向 Redis 提示哪些 keys，优先被淘汰。

**淘汰进程是如何工作的**

- client 运行新的命令，导致了添加更多数据。
- Redis 检测内存使用情况，如果大于 maxmemory limit，则根据策略淘汰。
- 运行新的命令，等等。



### 近似的 LRU 算法

```
maxmemory-samples 5
```

这就是为什么 Redis 不使用真正的 LRU 实现，因为这样会导致更多的内存使用。



### LFU mode (Least Frequently Used eviction mode)





