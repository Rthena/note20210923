## Redis 的分布式锁实现

### 单实例的正确实现方式

用如下，方法获取锁

```bash
    SET resource_name my_random_value NX PX 30000
```
该命令仅在 key不存在（NX 选项）时设置密钥，过期时间为 30000 毫秒（PX 选项）。 key 设置为值“my_random_value”。 

my_random_value必须要具有唯一性，我们可以用UUID来做，设置随机字符串保证唯一性，至于为什么要保证唯一性？假如value不是随机字符串，而是一个固定值，那么就可能存在下面的问题：

- 1.客户端1获取锁成功
- 2.客户端1在某个操作上阻塞了太长时间
- 3.设置的key过期了，锁自动释放了
- 4.客户端2获取到了对应同一个资源的锁
- 5.客户端1从阻塞中恢复过来，因为value值一样，所以执行释放锁操作时就会释放掉客户端2持有的锁，这样就会造成问题

1. 获取锁方式一，setnx命令

   ```go
   func TryLockWithSetNX(key, uniqueID string, timeout time.Duration) (bool, error) {
   	result, err := redisIns.SetNX(ctx, key, uniqueID, timeout).Result()
   	if err != nil {
   		fmt.Printf("%s lock err:%s", key, err)
   		return false, fmt.Errorf("%s lock err:%s", key, err)
   	}
   	if result {
   		fmt.Printf("%s lock success", key)
   		return true, nil
   	}
   	return false, fmt.Errorf("%s has been locked", key)
   }
   ```

   

2. 获取锁方式二，使用 lua 脚本

   ```go
   const lockScrip = "if redis.call('setnx',KEYS[1],ARGV[1]) == 1" +
   	" then redis.call('expire',KEYS[1],ARGV[2]) return 1 else return 0 end"
   
   func TryLockWithLua(key, uniqueID string, second int) (bool, error) {
   	result, err := redisIns.Eval(ctx, lockScrip, []string{key}, uniqueID, second).Result()
   	if err != nil {
   		fmt.Printf("%s lock err:%s", key, err)
   		return false, fmt.Errorf("%s lock err:%s", key, err)
   	}
   	if val, ok := result.(int64); ok && val == 1 {
   		fmt.Printf("%s lock success", key)
   		return true, nil
   	}
   	return false, fmt.Errorf("%s has been locked", key)
   }
   ```



解锁
```go
const releaseLockScrip = "if redis.call('get',KEYS[1]) == ARGV[1] then " +
	"return redis.call('del',KEYS[1]) else return 0 end"

func ReleaseLockWithLua(key, uniqueID string) (bool, error) {
	result, err := redisIns.Eval(ctx, releaseLockScrip, []string{key}, uniqueID).Result()
	if err != nil {
		fmt.Printf("%s release err:%s", key, err.Error())
		return false, fmt.Errorf("%s ReleaseLockWithLua err:%s", key, err.Error())
	}
	if val, ok := result.(int64); ok && val == 1 {
		fmt.Printf("%s release success", key)
		return true, nil
	}
	return false, fmt.Errorf("%s has been released", key)
}
```



### **为什么基于故障转移的实现是不够的**

表面上看，很有效，但有一个问题：这是一个单点故障架构，如果 Redis 主机宕机怎么办？。好吧， 让我们来添加一个副本，当主机不可用的时候。

这个模型有一个明显的竞争条件：
1. 客户端 A 从master获取锁。
2. master 把 key 传送到 replica 之前崩溃了
3. replica 晋升为 master
4. 客户端 B 获取了客户端 A已经持有同一资源的锁，违反了安全！



## Redlock算法

假设有5个独立的Redis节点（**注意这里的节点可以是5个Redis单master实例，也可以是5个Redis Cluster集群，但并不是有5个主节点的cluster集群**）：

- 获取当前Unix时间，以毫秒为单位
- 依次尝试从5个实例，使用相同的key和具有唯一性的value(例如UUID)获取锁，当向Redis请求获取锁时，客户端应该设置一个网络连接和响应超时时间，这个超时时间应用小于锁的失效时间，例如你的锁自动失效时间为10s，则超时时间应该在5~50毫秒之间，这样可以避免服务器端Redis已经挂掉的情况下，客户端还在死死地等待响应结果。如果服务端没有在规定时间内响应，客户端应该尽快尝试去另外一个Redis实例请求获取锁
- 客户端使用当前时间减去开始获取锁时间（步骤1记录的时间）就得到获取锁使用的时间，当且仅当从大多数(N/2+1，这里是3个节点)的Redis节点都取到锁，并且使用的时间小于锁失败时间时，锁才算获取成功。
- 如果取到了锁，key的真正有效时间等于有效时间减去获取锁所使用的时间（步骤3计算的结果）
- 如果某些原因，获取锁失败（没有在至少N/2+1个Redis实例取到锁或者取锁时间已经超过了有效时间），客户端应该在所有的Redis实例上进行解锁（即便某些Redis实例根本就没有加锁成功，防止某些节点获取到锁但是客户端没有得到响应而导致接下来的一段时间不能被重新获取锁）