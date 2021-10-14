## Redis 的分布式锁实现

1. **使用 setnx 命令（错误做法）**

   ```go
   func tryLock(key, reqSet string, timeout time.duration) bool {
       ret, err := redisIns.SetNX(ctx, key, reqSet, timeount).Result()
       if err != nil  {
           return false
       }
       if ret {
           return true
       }
       return false
   }
   ```

   这个模型有一个明显的竞争条件：

   1. 客户端 A 从master获取锁。
   2. master 把 key 传送到 replica 之前崩溃了
   3. replica 晋升为 master
   4. 客户端 B 获取了客户端 A已经持有同一资源的锁，违反了安全！

2. 

