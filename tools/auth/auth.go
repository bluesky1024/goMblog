package authen

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

//检查调用方是否有调用该服务的权限
func CheckPower(appId string, serverPath string) (token string, err error) {
	return "abc", err
}

//检查调用方签名
func CheckSign(sign string, token string, reqData interface{}) (err error) {
	//req数据获取
	mapData := make(map[string]string)
	v := reflect.ValueOf(reqData).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if ind := strings.Index(v.Type().Field(i).Name, "XXX_"); ind != 0 {
			mapData[v.Type().Field(i).Name] = fmt.Sprintf("%v", f)
		}
	}

	checkSign := GetSign(reqData, token)
	if checkSign != sign {
		return errors.New("auth fail, wrong sign")
	}
	return nil
}

func GetSign(data interface{}, token string) (sign string) {
	//req数据获取
	mapData := make(map[string]string)
	v := reflect.ValueOf(data).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		a := fmt.Sprintf("%v", f)

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
