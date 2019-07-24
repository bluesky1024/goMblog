package chatService

import (
	"github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func (s *chatService) StartRoom(uid int64) (err error) {
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

func (s *chatService) StopRoom(uid int64) error {
	return nil
}
