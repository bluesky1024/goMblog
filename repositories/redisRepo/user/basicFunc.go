package userRdRepo

import (
	"fmt"
)

//目前key都是当前包临时指定的名字
//后期redis的key都需要统一配置，将前缀替换为projectId_serviceId_FuncId_%d...的形式
//在保证不重复的前提下，尽量缩短key的长度
var (
	fanCntKey    = "FanCnt_%d"    //FanCtn_uid
	followCntKey = "FollowCnt_%d" //FollowCnt_uid
)

func getFanCntKey(uid int64) string {
	return fmt.Sprintf(fanCntKey, uid)
}

func getFollowCntKey(uid int64) string {
	return fmt.Sprintf(followCntKey, uid)
}
