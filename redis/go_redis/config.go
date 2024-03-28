package go_redis

import (
	"fmt"
	"time"
)

type Config struct {
	Host         string        `yaml:"host" json:"host,omitempty"`
	Port         int32         `yaml:"port" json:"port,omitempty"`
	Username     string        `yaml:"username" json:"username,omitempty"`
	Password     string        `yaml:"password" json:"password,omitempty"`
	MaxRetries   int           `json:"maxRetries,omitempty"`
	DialTimeout  time.Duration `json:"dialTimeout,omitempty"`
	ReadTimeout  time.Duration `json:"readTimeout,omitempty"`
	WriteTimeout time.Duration `json:"writeTimeout,omitempty"`
	PoolSize     int32         `json:"poolSize,omitempty"`
	MinIdleConns int           `json:"minIdleConns,omitempty"`
}

func (r *Config) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
