/**
 * @Author raven
 * @Description
 * @Date 2023/2/9
 **/
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

/** ToString
 * @Description: 将对象转换为string
 * @param value 要转换的对象，支持值和指针类型
 * @return string 返回的字符串值
 */
func ToString(value interface{}) (string, error) {
	ev := reflect.ValueOf(value)
	if reflect.ValueOf(value).Kind() == reflect.Ptr {
		ev = reflect.ValueOf(value).Elem()
	}
	if ev.Kind() == reflect.Struct || ev.Kind() == reflect.Map || ev.Kind() == reflect.Slice {
		// []byte类型
		bs, ok := ev.Interface().([]byte)
		if ok {
			return string(bs), nil
		}
		str, err := json.Marshal(ev.Interface())
		if err != nil {
			return "", errors.New(fmt.Sprintf("toString json marshal err:%v", err))
		}

		return string(str), nil
	} else {
		// string类型
		str, ok := ev.Interface().(string)
		if ok {
			return str, nil
		}
		// []byte类型
		bs, ok := ev.Interface().([]byte)
		if ok {
			return string(bs), nil
		}
		// int类型
		intVal, ok := ev.Interface().(int)
		if ok {
			return strconv.Itoa(intVal), nil
		}
	}

	return "", errors.New(fmt.Sprintf("toString unimplemented type:%v", reflect.TypeOf(value).Name()))
}
