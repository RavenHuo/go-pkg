/**
 * @Author raven
 * @Description
 * @Date 2021/8/31
 **/
package mongo

import (
	"errors"

	"gopkg.in/mgo.v2"
)

var (
	ErrMaxConnectionLimited = errors.New("mongo: reach max connections")
)

func IsErrMaxConnectionLimited(err error) bool {
	return err == ErrMaxConnectionLimited
}

func IsErrNotFound(err error) bool {
	return err == mgo.ErrNotFound
}
