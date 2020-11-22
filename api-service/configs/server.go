package config

import (
	"time"
)

// ServerConfiguration represent configs for gin-server
type ServerConfiguration struct {
	RunMode      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
