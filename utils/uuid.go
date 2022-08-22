/**
 * @Author raven
 * @Description
 * @Date 2022/7/25
 **/
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func GetUuid() string {
	buf := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, buf)
	return hex.EncodeToString(buf)
}
