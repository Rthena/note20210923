package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
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

func tryLock(key, reqSet string, timeout time.Duration) bool {
	result, err := redisIns.SetNX(ctx, key, reqSet, timeout).Result()
	if err != nil {
		fmt.Println("SetNX", err)
		return false
	}
	if result {
		return true
	}

	defer func() {

	}()
	return false
}

func TestTryLock(t *testing.T) {
	lock := tryLock("keynx", "666", 60*time.Second)
	t.Log(lock)
}

const luaScrip = "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"

func releaseLock(key string) {
	eval := redisIns.Eval(ctx, luaScrip, []string{key})
	result, err := eval.Result()
}
