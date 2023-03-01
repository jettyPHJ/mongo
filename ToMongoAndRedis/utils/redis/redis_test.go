package redis_test

import (
	"ToMongoAndRedis/utils/log"
	"testing"

	"github.com/go-redis/redis"
)

func TestRedisConnect(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.163:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := rdb.Set("mykey", "myvalue:jetty", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := rdb.Get("mykey").Result()
	if err != nil {
		panic(err)
	}
	log.ErrorLogger.Println("key", val)
	val2, err := rdb.Get("key2").Result()
	if err == redis.Nil {
		log.ErrorLogger.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		log.ErrorLogger.Println("key2", val2)
	}
}
