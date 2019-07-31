package chatRdRepo

import (
	"fmt"
	"github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/datasource/redisSource"
	"testing"
	"time"
)

var chatRd *ChatRbRepository

func init() {
	chatSource, err := redisSource.LoadChatRdSour()
	if err != nil {
		panic(err.Error())
	}
	chatRd = NewChatRdRepo(chatSource)
}

func TestChatRbRepository_GetRoomConfig(t *testing.T) {
	res, err := chatRd.GetRoomConfig(123)
	fmt.Println(res, err)
}

func TestChatRbRepository_SetRoomConfig(t *testing.T) {
	var uid int64 = 123
	config := datamodels.ChatRoomConfigure{
		RoomOwnerUid: uid,
		RedisSetCnt:  3,
		Status:       datamodels.RoomStatusNormal,
		WorkStatus:   datamodels.WorkStatusOn,
	}
	err := chatRd.SetRoomConfig(uid, config)
	fmt.Println(err)
}

func TestChatRbRepository_DelRoomConfig(t *testing.T) {
	var uid int64 = 123
	succ, err := chatRd.DelRoomConfig(uid)
	fmt.Println(succ, err)
}

func TestChatRbRepository_AppendNewBarrage(t *testing.T) {
	info := datamodels.ChatBarrageInfo{
		Uid:        111,
		Message:    "this is a test message",
		VideoTime:  22,
		CreateTime: time.Now(),
	}
	var uid int64 = 123
	err := chatRd.AppendNewBarrage(uid, 1, info)
	fmt.Println(err)
}

func TestChatRbRepository_GetBarragesByUid(t *testing.T) {
	var uid int64 = 123
	timeAfter := time.Unix(10000000, 0)
	//timeAfter = time.Now()
	res, err := chatRd.GetBarragesByUid(uid, 1, 5, timeAfter)
	fmt.Println(res, err)
}

func TestChatRbRepository_RemoveBarrageByUid(t *testing.T) {
	var uid int64 = 123
	timeBefore := time.Now()
	err := chatRd.RemoveBarrageByUid(uid, 1, timeBefore)
	fmt.Println(err)
}
