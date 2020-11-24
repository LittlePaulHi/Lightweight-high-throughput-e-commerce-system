package config

import (
	"time"
)

// Configuration of api-service
type Configuration struct {
	Server ServerConfiguration
	Cache  CacheConfiguration
}

// ServerConfiguration represent configs for gin-server
type ServerConfiguration struct {
	RunMode      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type CacheConfiguration struct {
	Host     string
	DataBase int
	Expires  time.Duration
}
