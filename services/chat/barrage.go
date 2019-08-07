package chatService

import (
	"errors"
	"github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"time"
)

var ()

func getSetIndByUid(uid int64, setCnt int) (setInd int) {
	setInd = int(uid%int64(setCnt)) + 1
	return setInd
}

func (s *chatService) GetBarrageByRoomId(uid int64, roomId int64) (barrages []datamodels.ChatBarrageInfo, err error) {
	//获取room配置
	roomConfig, err := s.rdRepo.GetRoomConfigByRoomId(roomId)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}

	//根据redisSetCnt进行hash，选择特定的set
	setInd := getSetIndByUid(uid, roomConfig.RedisSetCnt)

	//获取弹幕
	//每隔3秒获取50条足够了吧
	timeAfterStamp := time.Unix(time.Now().Unix()-5, 0)
	barrages, err = s.rdRepo.GetBarragesByUid(roomConfig.RoomOwnerUid, setInd, 50, timeAfterStamp)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return barrages, err
}

func (s *chatService) SendBarrage(uid int64, roomId int64, message string, videoTime int64) error {
	msg := datamodels.ChatBarrageInfo{
		Uid:        uid,
		RoomId:     roomId,
		Message:    message,
		VideoTime:  videoTime,
		CreateTime: time.Now(),
	}
	err := s.sendNewBarrageToRoom(msg)
	if err != nil {
		logger.Err(logType, err.Error())
		return errors.New("send barrage failed")
	}
	return nil
}
