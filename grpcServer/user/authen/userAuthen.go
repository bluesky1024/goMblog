package authen

import (
	"errors"
	"fmt"
	"github.com/bluesky1024/goMblog/tools/auth"
	"reflect"
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
			mapData[v.Type().Field(i).Name] = fmt.Sprintf("%v",f)
		}
	}

	checkSign := auth.GetSign(reqData, token)
	if checkSign != sign {
		return errors.New("auth fail, wrong sign")
	}
	return nil
}
