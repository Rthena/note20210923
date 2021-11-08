### [Redis 持久化](https://redis.io/topics/persistence)

Redis 提供了一系列不同的持久性选项：

- **RDB** (Redis Database)：以指定的时间间隔，执行数据集的时间点快照

- **AOF** (Append Only File)：服务器接收到的每个写操作，都会持久化到日志里。当redis服务器启动时会加载该日志文件重新把数据构建到内存中。命令使用 Redis协议本身相同的格式以 Append Only的方式记录的。当日志变得太大时，Redis能够在后台重写日志。

- **NO persistentce**：如果您只希望数据在服务器运行时才存在，那么您可以完全禁用持久性
  **RDB+AOF**：RDB，AOF可以同时使用。需要注意的是当Redis服务器重新启动的时候，默认是使用AOF 还原数据的，因为这样可以保证数据的完整性

  
  
  最重要的是你要理解它们之间有何区别，然后在权衡一下应该使用何种方式比较合适

**RDB 优点**

- RDB是一个非常紧凑的单文件。RDB文件非常合适于备份，例如，您希望在最近24小时内，每小时存档一次RDB文件，并且在30天内，每天保存一次快照。这使您可以在灾难时恢复不同的版本的数据。
- RDB非常有利于灾难恢复，单个压缩文件可传输到远程数据中心，也可以放到远程的服务器
- RDB 最大限度地提高了 Redis 的性能，因为 Redis 父进程为了持久化需要做的唯一工作是派生一个子进程将剩余工作完成。 父实例永远不会执行磁盘 I/O 或类似操作。
- 相比较于AOF 庞大的数据集，RDB启动得更快
- 备份方面，RDB支持 部分同步

**RDB 缺点**

- Redis服务没有正常退出时可能会丢失少量数据。每5分钟或者更长的时间创建一个快照，这个时候如果redis 没有正常退出，就会可能丢失最近几分钟内的数据。
- RDB经常需要fork()，以便子进程把内存数据保存到磁盘上，如果数据集很大，fork()会很耗时，可能会导致Redis停止对客户端提供服务几毫秒甚至1秒。如果dataset非常大，那么对CPU是不友好的。AOF也有fork()操作，但是你可以调整多久重写日志的频率，而不需要在耐用性上做任何权衡



**AOF 优点**

- 使用 AOF Redis会更持久。你可以有不同的 fsync 策略：no fsync at all, fsync every second, fsync at every query。默认的每秒钟同步策略，依然有很好的性能（fsync 是后台线程执行的，当没有fsync正进行时，主线程将尽可能的执行写入）。即使是数据就是，也只丢失 1秒中的写入。
- AOF 日志是一个只追加的日志，因此在断电之后不会出现寻道或者损坏问题。因此日志由于某种原因（磁盘已满或者其他原因）以半写命令结束，redis-check-aof 工具也能轻松修复。
- 当 AOF 变大时，Redis 能够在后台自动重写 AOF。 重写是完全安全的，因为当 Redis 继续追加到旧文件时，会使用创建当前数据集所需的最少操作集生成一个全新的文件，一旦第二个文件准备就绪，Redis 就会切换这两个文件并开始追加到 新的那一个。
- AOF 以易于理解和解析的格式包含所有操作的日志。 您甚至可以轻松导出 AOF 文件。 例如，即使您不小心使用 FLUSHALL 命令刷新了所有内容，只要在此期间没有重写日志，您仍然可以通过停止服务器、删除最新命令并重新启动 Redis 来保存您的数据集 再次。

**AOF 缺点**

- AOF 文件通常比相同数据集的等效 RDB 文件大。
- AOF 比RDB慢取决于采用那种 fsync 策略，通常来说将 fsync 设置成每秒，性能依然是很高的，并且禁用 fsync 后 ，即使在高负载的情况下AOF的也应该比RDB快。然而即使在写入负载巨大的情况下，RDB 仍然能够提供更多关于最大延迟的保证。
- 3、在过去的实验中 AOF在特定的命令中遇到过罕见的bugs ，导致生成的 AOF 在重新加载时无法重现完全相同的数据



**该如何选择**

- 一般的指示是，如果您想要与 PostgreSQL 可以提供的数据安全程度相当的数据安全性，则应该同时使用这两种持久性方法。
- 如果您非常关心您的数据，但在发生灾难时仍然可以忍受几分钟的数据丢失，您可以简单地单独使用 RDB。
- 数据库备份和灾难恢复，RDB快照非常便于数据库备份，并且RDB恢复数据集的速度比AOF 快。

**这两种持久化方式的更多细节**

**Snapshotting**

​	默认情况下，Redis会将数据集快照保存在磁盘上一个名为 dump.rdb的二进制文件中。你可以设置 Redis在 N 秒内最少有 M 次修改保存一次数据集，或者您可以手动调用 SAVE 或 BGSAVE 命令

例如，如下配置，60秒内最少有1000个key 被修改，Redis 将自动把数据集保存到磁盘

```
save 60 1000
```

**工作原理**

- Redis forks，子进程和父进程。
- 子进程把数据集写入到一个临时文件中。
- 当子进程完成写临时文件后，它会替换老的 RDB文件

这样的好处是copy-on-write



**Append-only file**

快照方式 *Snapshotting*，并不是那么的可靠，如果你运行Redis的计算机停止工作了，你的电源线坏了，或者是你不小心执行了 kill -9 命令了，这就意味着最近写入到 Redis的数据丢失了。对于一些应用来说这并不是什么大问题，但是对于需要完全可靠的场景来说，Redis不是可行的选择。

对于完全可靠场景 *Append-only file* 模式非常合适。

您可以在配置文件中打开 AOF

```
appendonly yes
```



***Log rewriting 日志重写***

正如你可以猜想到的，当写操作被执行时 AOF 会越来越大。例如，对于一个计数器，加100次，您最终数据集中只有一个 key 包含最终值，但是这个时候 AOF 中有 100 个条目，重建当前状态不需要这些条目中的 99 个。

因此 Redis 支持一个很有趣但特性，它可以在后台重建 AOF，无需打断对 client 的服务。当你发出 BGREWRITEAOF 时，Redis 都会写入内存中重建当前数据集所需的最短命令序列。如果您在 Redis 2.2 中使用 AOF，则需要时不时运行 BGREWRITEAOF。Redis 2.4 能自动触发日志重新 

***Append only file* 是怎样持久化的？**

您可以配置 Redis 将在磁盘上同步数据的次数。 共有三个选项：

- appendfsync always:  每次有数据修改发生时都会写入AOF文件。
- appendfsync everysec: 每秒钟保存一次，足够快，即使是发生灾难，也只会丢失 1秒钟的数据
- appendfsync no: 不进行同步，只是把你的数据交给操作系统，速度相当的快，却并不安全。通常 Linux 每30秒 flush 数据，但是这只决定于内核是怎样调节的

默认策略是 fsync everysec，速度快并且相当的安全 ，always 策略在实践中很慢，但它支持组提交 group commit，所以如果有多个并行写入 Redis 可以尝试执行单个 fsync 操作。



**AOF 被截断了，应该处理？**

这个情况下 Redis会抛出如下日志

```
* Reading RDB preamble from AOF file...
* Reading the remaining AOF tail...
# !!! Warning: short read while loading the AOF file !!!
# !!! Truncating the AOF at offset 439 !!!
# AOF loaded anyway because aof-load-truncated is enabled
```

发生这种情况，有可能在服务器写入 AOF 文件的时崩溃了，或者写入 AOF 文件所在的卷 volume已满。 在这种情况时，AOF中最后一个命令可能会被截断。最新版的 Redis无论如何都可以加载 AOF，只需要丢弃文件中最后一个格式不正确的命令。

对于旧版 Redis，按照如下步骤操作

- 备份一份 AOF文件
- 使用 redis-check-aof 工具修复 原来的 AOF文件 运行命令 $ redis-check-aof --fix
- 可选，使用 diff -u 检测 两个文件之间的差异
- 使用修复后的文件，重启 Redis

**AOF 文件损坏了，怎样处理？**

AOF文件没有被截断，可是文件中的无效字节顺序损坏了，事情会变得更加复杂，Redis会抛出启动失败日志

```
* Reading the remaining AOF tail...
# Bad file format reading the append only file: make a backup of your AOF file, then use ./redis-check-aof --fix <filename>
```

使用命令 $ redis-check-aof --fix，了解问题，跳转到，获取给定文件的偏移量，尝试手动修复文件。否则有可能让实用程序为我们修复文件，但在这种情况下，从无效部分到文件末尾的所有 AOF 部分可能会被丢弃，如果损坏发生，将导致大量数据丢失

**工作原理 （copy-on-write）**

- Redis forks，因此现在我们有了父子进程。
- 子进程，开始写新的 AOF临时文件。
- 父进程累积所有的new changes 到内存缓冲区中（于此同时，父进程也会往 old AOF文件中写入 new changes，因此如果重写失败了，数据还是安全的）。
- 当子进程重写文件完毕，父进程会收到一个信号，并将内存缓冲区追加到子进程生成的文件的末尾。
- 完美！ 现在 Redis 原子地将旧文件重命名为新文件，并开始将新数据追加到新文件中。

**当前使用的 dump.rdb 快照方式，怎样切换到 AOF？**

**Redis >= 2.2** （无需重启）

- 备份最新的 dump.rdb 文件
- 把备份文件转移到安全的地方
- 运行以下的两条命令
- redis-cli config set appendonly yes
- redis-cli config set save ""
- 确保您的数据库包含与AOF文件包含的相同数量的键。
- 确保追加的文件正确。

第一条命令是用来启用 Append Only File的，为此Reids将阻塞 dump 初始化，然后将打开文件进行写入，并开始追加下一个写入查询。

第二条命令是用来关闭 snapshotting persistence的。

<font size=5>重要</font> ：记得要修改 redis.conf 文件打开 AOF ！！，否则的话，当您重启服务器的时候修改的配置都会丢失，服务器将会加载 old 配置

**Redis 2.0**

- 备份最新的 dump.rdb 文件
- 把备份文件转移到安全的地方
- 停止对数据库的所有写操作！
- 运行以下的两条命令
- redis-cli config set appendonly yes
- redis-cli config set save ""
- 确保您的数据库包含与AOF文件包含的相同数量的键。
- 确保追加的文件正确。

第一条命令是用来启用 Append Only File的，为此Reids将阻塞 dump 初始化，然后将打开文件进行写入，并开始追加下一个写入查询。

第二条命令是用来关闭 snapshotting persistence的。

<font size=5>重要</font> ：记得要修改 redis.conf 文件打开 AOF ！！，否则的话，当您重启服务器的时候修改的配置都会丢失，服务器将会加载 old 配置



### AOF 和 RDB 持久化之间的交互

Reids >= 2.4 当一个 RDB 快照正在执行的时候 Redis会避免触发 AOF重写，当 AOF 重写正在执行时 Redis 允许使用 BGSAVE命令。这样可以避免两个后台进程同时进行很重的磁盘 I/O操作。当 snapshotting 进程正在执行时，用户使用 BGREWRITEAOF 命令明确的发起日志重写操作，服务器将会返回 OK，告诉用户操作已经被安排好，一旦 snapshotting 执行完毕，重写日志将会开始。

在 AOF 和 RDB 持久化都启用并且 Redis 重启的情况下，AOF 文件将用于重建原始数据集，因为它保证是最完整的。



### 备份 Redis数据

Redis数据备份是非常友好的，尽管数据库还在运行，您也可以拷贝RDB文件，RDB文件一旦生成就不会再修改了，它在产生时会使用一个临时名字，当新的快照完成时，它会自动的重命名最终的名字。

这就意味着即使服务器还在运行，拷贝 RDB文件是完全安全的，这里有几点建议

- 在您的服务器，创建一个 cron job ，在一个目录中创建 RDB文件的每小时快照，并在另一个目录中创建每日快照。
- 每次执行这个脚本，确保调用查找命令，确保太老 快照被删除：例如您可以拿最近 48小时的 每小时快照，以及一两个月的每日快照。确保快照名字都有日期时间信息
- 确保每天至少一次转移 RDB快照 到数据中心之外 或者 Redis实例的物理机之外

如果您在仅启用 AOF 持久性的情况下运行 Redis 实例，您仍然可以复制 AOF 以创建备份。 该文件可能缺少最后一部分，但 Redis 仍然可以加载它



### 灾难恢复

Redis上下文中的灾难恢复于备份基本相同，并且能够中许多不同的外部数据中心传输这些备份。通过这种方式，即使在影响 Redis运行并生成其快照的主数据中心发生灾难性事件，数据也能得到保护。

- 传输到远程服务器，每天或者每小时传递带有加密的快照到远程服务器。把密码保存到多个地方，例如，可以给组织中重要的人。
- 通过SCP 传输到远程服务器。

