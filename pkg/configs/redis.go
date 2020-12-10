package config

import (
	"time"
)

type RedisConfiguration struct {
	Address         string
	DataBase        int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolSize        int
	PoolTimeout     time.Duration
	CacheExpireTime time.Duration
}
