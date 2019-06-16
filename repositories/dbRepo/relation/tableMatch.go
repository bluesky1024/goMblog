package relationDbRepo

import (
	"strconv"
)

var followTableName = "follow_info"
var fanTableName = "fan_info"
var followGroupTableName = "follow_group"
var relationCntTableName = "relation_cnt"

// 根据mid中时间信息获取follow的表名
func getFollowTableName(uid int64) string {
	virtualNode := uid % 10000
	tableNum := virtualNode/1000 + 1
	return followTableName + "_" + strconv.FormatInt(tableNum, 10)
}

// 根据uid中时间信息获取fan的表名
func getFanTableName(uid int64) string {
	virtualNode := uid % 10000
	tableNum := virtualNode/1000 + 1
	return fanTableName + "_" + strconv.FormatInt(tableNum, 10)
}

// 根据uid中时间信息获取follow_group的表名
func getFollowGroupTableName(uid int64) string {
	virtualNode := uid % 10000
	tableNum := virtualNode/1000 + 1
	return followGroupTableName + "_" + strconv.FormatInt(tableNum, 10)
}

// 根据uid中时间信息获取follow_group的表名
func getRelationCntTableName(uid int64) string {
	virtualNode := uid % 10000
	tableNum := virtualNode/1000 + 1
	return relationCntTableName + "_" + strconv.FormatInt(tableNum, 10)
}
