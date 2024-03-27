package go_redis

import (
	"fmt"
	"time"
)

type Config struct {
	Host         string
	Port         int32
	Username     string
	Password     string
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int32
	MinIdleConns int
}

func (r *Config) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
