package mblogDbRepo

import (
	//"fmt"
	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"
	//"strconv"
)

var baseMblogTableName string = "mblog_info"
var baseUidToMblogTableName string = "uid_to_mblog"

// 根据mid中时间信息获取mblog的表名
func getMblogTableName(mid int64) string {
	timeFormat := idGen.GetTimeInfoById(mid).Format("200601")
	return baseMblogTableName + "_" + timeFormat
}

// 根据uid中时间信息获取uid_to_mblog的表名(此处改成按uid一致性取模hash来获取表名)
func getUidToMblogTableName(uid int64) string {
	timeFormat := idGen.GetTimeInfoById(uid).Format("200601")
	return baseUidToMblogTableName + "_" + timeFormat
}

////虚拟节点10000个，目前初定10张表(0~999=>uid_to_mblog_1,1000~1999=>uid_to_mblog_2,...)
//func getUidToMblogTableName(uid int64) string {
//	virtualNode := uid % 10000
//	fmt.Println(virtualNode)
//	tableNum := virtualNode/1000 + 1
//	fmt.Println(tableNum)
//	return baseUidToMblogTableName + "_" + strconv.FormatInt(tableNum, 10)
//}
