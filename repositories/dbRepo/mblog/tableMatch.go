package mblogDbRepo

import (
	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"
)

var baseMblogTableName string = "mblog_info"
var baseUidToMblogTableName string = "uid_to_mblog"

// 根据mid中时间信息获取mblog的表名
func getMblogTableName(mid int64) string {
	timeFormat := idGen.GetTimeInfoById(mid).Format("200601")
	return baseMblogTableName + "_" + timeFormat
}

// 根据uid中时间信息获取uid_to_mblog的表名(此处改成按uid取模hash来获取表名)
func getUidToMblogTableName(uid int64) string {
	timeFormat := idGen.GetTimeInfoById(uid).Format("200601")
	return baseUidToMblogTableName + "_" + timeFormat
}
