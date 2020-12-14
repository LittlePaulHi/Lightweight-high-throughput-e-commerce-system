package config

import (
	"time"
)

type MariadbConfiguration struct {
	Type            string
	User            string
	Host            string
	Name            string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}
