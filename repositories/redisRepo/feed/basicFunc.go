package feedRdRepo

import (
	"fmt"
)

var (
	feedKey = "%d_%d" //uid_groupId
)

func getFeedKey(uid int64, groupId int64) string {
	return fmt.Sprintf(feedKey, uid, groupId)
}
