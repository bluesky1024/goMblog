package chatService

import (
	"github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func (s *chatService) AddRoom(roomName string, roomId int64, roomOwnerUid int64, redisSetCnt int) error {
	info := datamodels.ChatRoomConfigure{
		RoomName:     roomName,
		RoomId:       roomId,
		RoomOwnerUid: roomOwnerUid,
		RedisSetCnt:  redisSetCnt,
		Status:       datamodels.RoomStatusNormal,
	}
	err := s.dbRepo.AddRoom(info)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}

func (s *chatService) RemoveRoom(roomId int64) error {
	err := s.dbRepo.RemoveRoom(roomId)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}
