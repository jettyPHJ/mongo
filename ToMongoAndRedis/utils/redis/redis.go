package redis

import (
	cfg "ToMongoAndRedis/utils/config"

	"github.com/go-redis/redis"
)

var Client *redis.Client

func init() {

	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Param.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
