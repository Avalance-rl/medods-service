package cache

import (
	"medods-service/internal/config"

	"github.com/go-redis/redis"
)



var (
    RDB *redis.Client
)


func InitRedis() {
    RDB = redis.NewClient(&redis.Options{
        Addr: config.CFG.Redis.Address,
        Password: config.CFG.Redis.Password,
        DB:       config.CFG.Redis.DB, 
    })
}



