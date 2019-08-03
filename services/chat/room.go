package chatService

import (
	"errors"
	"github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"strconv"
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

func (s *chatService) GetRoomConfigByRoomId(roomId int64) (info datamodels.ChatRoomConfigure, err error) {
	info, err = s.rdRepo.GetRoomConfigByRoomId(roomId)
	if err != nil {
		logger.Err(logType, err.Error())
		var found bool
		info, found = s.dbRepo.GetRoomConfigByRoomId(roomId)
		if !found {
			err = errors.New(strconv.FormatInt(roomId, 10) + "room not found")
			logger.Err(logType, err.Error())
			return info, err
		}
		err = s.rdRepo.SetRoomConfig(roomId, info)
		if err != nil {
			logger.Err(logType, err.Error())
		}
	}
	return info, nil
}

func (s *chatService) StartRoomReal(uid int64) (err error) {
	defer func() {
		if err != nil {
			releaseErr := s.releaseRoomResource(uid)
			if releaseErr != nil {
				logger.Err(logType, releaseErr.Error())
			}
			changerErr := s.dbRepo.ChangeRoomWorkStatus(uid, datamodels.WorkStatusOn, datamodels.WorkStatusInvalid, false)
			if changerErr != nil {
				logger.Err(logType, changerErr.Error())
			}
		}
	}()
	//更改房间状态为开启中
	err = s.dbRepo.ChangeRoomWorkStatus(uid, datamodels.WorkStatusInvalid, datamodels.WorkStatusStarting, true)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}

	//初始化直播房间需要的资源
	err = s.initRoomResource(uid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}

	//更改房间状态为直播中
	err = s.dbRepo.ChangeRoomWorkStatus(uid, datamodels.WorkStatusStarting, datamodels.WorkStatusOn, true)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	return nil
}

//主播发送启动房间的指令，发送
func (s *chatService) StartRoom(uid int64) error {
	err := s.sendRoomStartMsg(uid)
	if err != nil {
		logger.Err(logType, err.Error())
		return errors.New(strconv.FormatInt(uid, 64) + "send start room msg failed")
	}
	return nil
}

func (s *chatService) StopRoom(uid int64) error {
	err := s.sendRoomStopMsg(uid)
	if err != nil {
		logger.Err(logType, err.Error())
		return errors.New(strconv.FormatInt(uid, 64) + "send stop room msg failed")
	}
	return nil
}
