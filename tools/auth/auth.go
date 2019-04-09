package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func GetSign(data interface{}, token string) (sign string) {
	//req数据获取
	mapData := make(map[string]string)
	v := reflect.ValueOf(data).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		a := fmt.Sprintf("%v",f)

		ind := strings.Index(v.Type().Field(i).Name, "XXX_")
		if ind != 0 {
			mapData[v.Type().Field(i).Name] = a
		}
	}

	sortArr := make([]string, len(mapData))
	aInd := 0
	for ind, _ := range mapData {
		sortArr[aInd] = ind
		aInd++
	}

	sort.Strings(sortArr)
	sortedStr := ""
	for _, v := range sortArr {
		sortedStr += mapData[v]
	}
	sortedStr += token

	//md5值
	h := md5.New()
	h.Write([]byte(sortedStr))
	return hex.EncodeToString(h.Sum(nil))
}
