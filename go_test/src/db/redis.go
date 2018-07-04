package dbMange

import (
	"github.com/go-redis/redis"
	"fmt"
)

func ExampleNewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:9999",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}

func ExampleClient(client *redis.Client) {
	err := client.Set("lin", "hello", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("lin").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("lin", val)

	val2, err := client.Get("lin2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: lin hello
	// lin2 does not exist
}

// 测试使用redis
func TestConnRedis() {
	client := ExampleNewClient()
	ExampleClient(client)
}