package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	_ "time"
)

var ctx = context.Background()

//var keyschan = make(chan string, 1024)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var wg = sync.WaitGroup{}

func ExampleClient() {
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	keys, _, err := rdb.Scan(ctx, 0, "*", 1000000).Result()

	if err != nil {
		panic(err)
	}

	for _, val := range keys {
		fmt.Printf("get key  %s\n", val)
		wg.Add(1)
		go getTtl(val)
	}

	wg.Wait()

}

func main() {
	ExampleClient()
}

func getTtl(key string) {
	key_type, err := rdb.Type(ctx, key).Result()

	if err != nil {
		panic(err)
	}

	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	switch ttl {
	case -1:
		if key_type == "string" {
			fmt.Printf("%s %s has no ttl\n", key_type, key)
			idel_time, err := rdb.ObjectIdleTime(ctx, key).Result()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s idel time is %d\n", key, idel_time/1000/1000/1000/60/60/24)
		}
	default:
		//fmt.Printf("%s ttl %d \n", key, ttl)
	}
	wg.Done()

}
