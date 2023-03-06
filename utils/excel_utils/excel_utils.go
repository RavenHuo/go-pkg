/**
 * @Author raven
 * @Description
 * @Date 2023/3/5
 **/
package excel_utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
)

const rowTagName = "row"

var UnSupportTypeError = errors.New("unsupport type")

// ReadExcel 读取excel内容
func ReadExcel(fileName string) ([][]string, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Printf("open excel error: %v\n", err)
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("close excel error: %v \n", err)
		}
	}()
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Printf("get excel rows error: %v\n", err)
		return nil, err
	}
	return rows, nil
}

// ParseExcel 读取excel数据，通过row tag 转换成对象
// @param fileName    文件名
// @param resultPtr  是需要解析的数组interface，只接受数组或者slice
// @param breakHeader 是否跳过第一行表头
func ParseExcel(fileName string, resultPtr interface{}, breakHeader bool) error {
	excelData, err := ReadExcel(fileName)
	if err != nil {
		return err
	}
	var resValue = reflect.ValueOf(resultPtr)
	// 只能是slice
	if resValue.IsNil() || resValue.Kind() != reflect.Ptr {
		return UnSupportTypeError
	}
	// 获取 resultPtr 里面的类型
	slicev := resValue.Elem()
	if slicev.Kind() != reflect.Array && slicev.Kind() != reflect.Slice {
		return UnSupportTypeError
	}
	// 根据 slicev 创建一个空slice
	slicev = slicev.Slice(0, slicev.Cap())
	//fmt.Println("input slice type of ", slicev.Type())
	// slice里面的类型
	elemt := slicev.Type().Elem()

	// slice里面类型也必须是指针
	if elemt.Kind() != reflect.Ptr {
		return UnSupportTypeError
	}
	// 获取 slice里面类型的值
	sliceType := elemt.Elem()

	fmt.Println("slice inner elem type ", sliceType.Kind())

	for i := 0; i < len(excelData); i++ {
		// 第一行
		if i == 0 && breakHeader {
			continue
		}

		// 利用反射创建一个新对象
		elemI := reflect.New(sliceType).Interface()

		err = setResult(excelData[i], elemI)
		if err != nil {
			return err
		}
		slicev = reflect.Append(slicev, reflect.ValueOf(elemI))
		slicev = slicev.Slice(0, slicev.Cap())
	}

	resValue.Elem().Set(slicev)
	return nil
}

func setResult(excelRow []string, resultPtr interface{}) error {
	if len(excelRow) == 0 {
		return nil
	}

	//
	var resType = reflect.TypeOf(resultPtr)

	// 判断类型
	if resType == nil || resType.Kind() != reflect.Ptr {
		return UnSupportTypeError
	}
	// 判断值
	var resElem = reflect.TypeOf(resultPtr).Elem()
	if resElem.Kind() != reflect.Struct {
		return UnSupportTypeError
	}

	for i := 0; i < resElem.NumField(); i++ {
		field := resElem.Field(i)
		rowTag := field.Tag.Get(rowTagName)

		rowIndex, convertErr := strconv.Atoi(rowTag)
		if convertErr != nil {
			continue
		}
		// 避免下标溢出
		if len(excelRow) <= rowIndex {
			continue
		}
		// excel的内容
		rowData := excelRow[rowIndex]
		var rowVal interface{}
		var err error
		fieldType := field.Type
		// 根据类型转换 rowVal
		if fieldType.Kind() == reflect.Struct || fieldType.Kind() == reflect.Map ||
			fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array {
			rowVal = reflect.New(fieldType).Interface()
			err = json.Unmarshal([]byte(rowData), rowVal)
			if err != nil {
				return err
			}

		} else if fieldType.Kind() == reflect.Int { //[]byte转换到int
			rowVal, err = strconv.Atoi(rowData)
			if err != nil {
				return err
			}
		} else if fieldType.Kind() == reflect.String { // []byte转换到string
			rowVal = rowData
		} else {
			return UnSupportTypeError
		}
		// resultVal被修改需要重新设置valType
		valType := reflect.TypeOf(rowVal)

		// 再次检查是否可设置值
		if valType.AssignableTo(fieldType) {
			// 指针不能设置值，只能用类型设置值
			reflect.ValueOf(resultPtr).Elem().Field(i).Set(reflect.ValueOf(rowVal))
		}
	}

	return nil
}
