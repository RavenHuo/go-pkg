/**
 * @Author raven
 * @Description
 * @Date 2023/2/9
 **/
package utils

import (
	"context"
	"encoding/json"
	"github.com/RavenHuo/go-kit/log"
	"reflect"
	"strconv"
)

/** ToString
 * @Description: 将对象转换为string
 * @param value 要转换的对象，支持值和指针类型
 * @return string 返回的字符串值
 */
func ToString(value interface{}) string {
	ev := reflect.ValueOf(value)
	if reflect.ValueOf(value).Kind() == reflect.Ptr {
		ev = reflect.ValueOf(value).Elem()
	}
	if ev.Kind() == reflect.Struct || ev.Kind() == reflect.Map || ev.Kind() == reflect.Slice {
		// []byte类型
		bs, ok := ev.Interface().([]byte)
		if ok {
			return string(bs)
		}
		str, err := json.Marshal(ev.Interface())
		if err != nil {
			log.Errorf(context.Background(), "toString json marshal err:%v", err)
		}

		return string(str)
	} else {
		// string类型
		str, ok := ev.Interface().(string)
		if ok {
			return str
		}
		// []byte类型
		bs, ok := ev.Interface().([]byte)
		if ok {
			return string(bs)
		}
		// int类型
		intVal, ok := ev.Interface().(int)
		if ok {
			return strconv.Itoa(intVal)
		}
	}

	log.Errorf(context.Background(), "toString unimplemented type:%v", reflect.TypeOf(value).Name())
	return ""
}
