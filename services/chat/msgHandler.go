package chatService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func (s *chatService) HandleRoomStartMsg(msg dm.RoomStatusSwitchMsg) (err error) {
	defer func() {
		if err != nil {
			releaseErr := s.releaseRoomResource(msg.Uid)
			if releaseErr != nil {
				logger.Err(logType, releaseErr.Error())
			}
			changerErr := s.dbRepo.ChangeRoomWorkStatus(msg.Uid, dm.WorkStatusOn, dm.WorkStatusInvalid, false)
			if changerErr != nil {
				logger.Err(logType, changerErr.Error())
			}
		}
	}()
	//更改房间状态为开启中
	err = s.dbRepo.ChangeRoomWorkStatus(msg.Uid, dm.WorkStatusInvalid, dm.WorkStatusStarting, true)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}

	//初始化直播房间需要的资源
	err = s.initRoomResource(msg.Uid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}

	//更改房间状态为直播中
	err = s.dbRepo.ChangeRoomWorkStatus(msg.Uid, dm.WorkStatusStarting, dm.WorkStatusOn, true)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	return nil
}

func (s *chatService) HandleRoomStopMsg(msg dm.RoomStatusSwitchMsg) (err error) {
	//更改房间状态为开始停止直播
	err = s.dbRepo.ChangeRoomWorkStatus(msg.Uid, dm.WorkStatusOn, dm.WorkStatusStoping, true)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}

	//释放直播房间需要的资源
	err = s.releaseRoomResource(msg.Uid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}

	//更改房间状态为停播中
	err = s.dbRepo.ChangeRoomWorkStatus(msg.Uid, dm.WorkStatusStoping, dm.WorkStatusInvalid, true)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	return nil
}

func (c *chatService) HandleNewBarrageToRoomMsg(msg dm.ChatBarrageInfo) (err error) {

	return nil
}
