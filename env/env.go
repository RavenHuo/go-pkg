package env

import (
	"os"
	"strings"
)

const (
	DEV_ENV  = "dev"  //开发环境
	PRO_ENV  = "pro"  //线上环境
	TEST_ENV = "test" //测试环境
)

var env string

func init() {
	env = os.Getenv("run_mode")
	env = strings.ToLower(env)
}

func GetEnv() string {
	return env
}
