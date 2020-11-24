package api

import (
	"api-service/cache"
	"api-service/config"
)

var redisCartCache cache.CartCache

func InitAllCache(conf *config.CacheConfiguration) {

	redisCartCache = cache.NewRedisCartCache(conf.Host, conf.DataBase, conf.Expires)
}
