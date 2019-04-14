package userGrpc
//package main
//
//import(
//	"fmt"
//	pb "github.com/bluesky1024/goMblog/services/userGrpc/userProto"
//	"github.com/golang/protobuf/ptypes/timestamp"
//	//"strconv"
//	"time"
//	"reflect"
//	"strings"
//)
//
//func main() {
//	uTime := timestamp.Timestamp{
//		Seconds:time.Now().Unix(),
//	}
//	a := pb.User{
//		Uid:123,
//		NickName:"abc",
//		CreateTime:&uTime,
//		UpdateTime:&uTime,
//	}
//	testMapString(&a)
//}
//
//func testMapString(a interface{}) {
//	//req数据获取
//	mapData := make(map[string]string)
//	v := reflect.ValueOf(a).Elem()
//
//	for i := 0; i < v.NumField(); i++ {
//		f := v.Field(i)
//		//fmt.Println(f)
//
//		a := fmt.Sprintf("%v",f)
//
//		indA := strings.Index(v.Type().Field(i).Name, "XXX")
//		//indB := strings.Index(v.Type().Field(i).Name, "time")
//		if indA != 0 {
//			mapData[v.Type().Field(i).Name] = a
//		}
//		//fmt.Println(v.Type().Field(i).Name)
//	}
//
//	fmt.Println(mapData)
//
//	//fmt.Println(reflect.TypeOf(a))
//}