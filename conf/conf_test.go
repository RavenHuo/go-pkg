package conf

import (
	"fmt"
	"testing"
)

type RegisterConfig struct {
	EtcdAddr     []string `json:"etcd_addr" yaml:"etcd_addr"`
	KeepaliveTtl int      `json:"keepalive_ttl" yaml:"keepalive_ttl"`
}

func TestGetObject(t *testing.T) {
	InitConf()
	registerConfig := new(RegisterConfig)
	GetObject("register", &registerConfig)
	fmt.Println(registerConfig.KeepaliveTtl)
}
func TestGetString(t *testing.T) {
	InitConf()
	fmt.Println(GetString("name"))
}
func TestGetInt(t *testing.T) {
	InitConf()
	fmt.Println(GetInt("server.name"))
}
