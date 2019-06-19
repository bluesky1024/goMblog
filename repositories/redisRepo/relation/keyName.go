package relationRdRepo

import (
	"fmt"
)

var (
	fanCntKey    = "fanCnt_%d"    //uid
	followCntKey = "followCnt_%d" //uid
)

func getFanCntKey(uid int64) string {
	return fmt.Sprintf(fanCntKey, uid)
}

func getFollowCntKey(uid int64) string {
	return fmt.Sprintf(followCntKey, uid)
}
