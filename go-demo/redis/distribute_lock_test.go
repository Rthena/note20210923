package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"testing"
	"time"
)

var redisIns *redis.Client

func init() {
	redisIns = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

const lock = "tryLock"

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
func TestTryLockWithSetNX(t *testing.T) {
	lock, err := TryLockWithSetNX(lock, "1", 60*time.Second)
	t.Log(lock, err)
}

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
func TestTryLockWithLua(t *testing.T) {
	lock, err := TryLockWithLua(lock, "1", 300)
	t.Log(lock, err)
}

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
func TestReleaseLockWithLua(t *testing.T) {
	lock, err := ReleaseLockWithLua(lock, "1")
	t.Log(lock, err)
}

func TestEval(t *testing.T) {
	result, err := redisIns.Eval(ctx, "return {KEYS[1],ARGV[1]}", []string{"key"}, "hello").Result()
	fmt.Println(result, err)
}

func TestRedisSync(t *testing.T) {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	mutexname := "my-global-mutex"
	mutex := rs.NewMutex(mutexname)

	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.
	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	// Do your work that requires the lock.
	fmt.Println("aaa")

	// Release the lock so other processes or threads can obtain a lock.
	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
}
