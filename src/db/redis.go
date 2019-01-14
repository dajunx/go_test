package dbMange

import (
	"github.com/go-redis/redis"
	"fmt"
)

func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		//Addr:     "localhost:9999",
		//Addr:     "192.168.221.137:6379",
		Addr:     "172.20.23.88:6379",
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

// Hash 数据结构相关
func SaveHashDatas(key string, input map[string]interface{}) {
	client := NewClient()
	err := client.HMSet(key, input).Err()
	if err != nil {
		panic(err)
	}
}

func GetHashDatas(key string) (map[string]string, error) {
	client := NewClient()
	val, err := client.HGetAll(key).Result()
	return val, err
}

// 测试使用redis
func TestConnRedis() {
	client := NewClient()
	ExampleClient(client)
}