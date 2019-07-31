package chatRdRepo

import (
	"fmt"
	"time"
)

var (
	//room config
	roomConfigKey          = "roomConfig_%d"
	roomConfigTimeDuration = 600 * time.Second

	//barrage set
	roomBarrageSetKey = "barrage_set_%d_%d"
)

func getRoomConfigSetInfo(uid int64) (key string, expireTime time.Duration) {
	return fmt.Sprintf(roomConfigKey, uid), roomConfigTimeDuration
}

func getBarrageSetInfo(uid int64, setInd int) (key string) {
	return fmt.Sprintf(roomBarrageSetKey, uid, setInd)
}
