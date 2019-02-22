package relationDbRepo

import (
	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"
)

var followTableName = "follow_info"
var fanTableName = "fan_info"
var followGroupTableName = "follow_group"

// 根据mid中时间信息获取follow的表名
func getFollowTableName(uid int64)  string{
	timeFormat := idGen.GetTimeInfoById(uid).Format("200601")
	return followTableName+"_"+timeFormat
}

// 根据uid中时间信息获取fan的表名
func getFanTableName(uid int64) string{
	timeFormat := idGen.GetTimeInfoById(uid).Format("200601")
	return fanTableName+"_"+timeFormat
}

// 根据uid中时间信息获取follow_group的表名
func getFollowGroupTableName(uid int64) string{
	timeFormat := idGen.GetTimeInfoById(uid).Format("200601")
	return followGroupTableName+"_"+timeFormat
}