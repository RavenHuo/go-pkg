package conf

import (
	"fmt"
	"testing"
)

func TestLoadConf(t *testing.T) {
	c := Init()
	fmt.Println(c.Get("name"))
}
