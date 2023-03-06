/**
 * @Author raven
 * @Description
 * @Date 2023/3/5
 **/
package excel_utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

type user struct {
	UserId int    `row:"0" json:"user_id"`
	Name   string `row:"1" json:"name"`
	Age    int    `row:"2" json:"age"`
}

func TestParseExcel(t *testing.T) {
	var result []*user
	err := ParseExcel("example.xlsx", &result, true)
	if err != nil {
		fmt.Printf("err :%s", err.Error())
	}
	jsonByte, _ := json.Marshal(result)
	fmt.Printf("r :%+v", string(jsonByte))
}
