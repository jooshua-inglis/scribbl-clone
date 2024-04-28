package eventListener

import (
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client = nil

func GetPubSub() *redis.Client {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}
	return rdb
}
