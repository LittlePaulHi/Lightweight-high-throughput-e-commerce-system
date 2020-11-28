package api

import (
	"api-service/cache"
	"api-service/config"
)

var (
	redisCartCache    cache.CartCache
	redisOrderCache   cache.OrderCache
	redisProductCache cache.ProductCache
)

func InitAllCache(conf *config.CacheConfiguration) {

	redisCartCache = cache.NewRedisCartCache(conf.Host, conf.DataBase, conf.Expires)
	redisOrderCache = cache.NewRedisOrderCache(conf.Host, conf.DataBase, conf.Expires)
	redisProductCache = cache.NewPostOrderCache(conf.Host, conf.DataBase, conf.Expires)
}
