package chatDbRepo

import (
	"errors"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"github.com/go-xorm/xorm"
)

var logType string = "chatDbRepo"

type ChatDbRepository struct {
	sourceM *xorm.Engine
	sourceS *xorm.Engine
}

func NewChatRepository(sourceM *xorm.Engine, sourceS *xorm.Engine) *ChatDbRepository {
	return &ChatDbRepository{
		sourceM: sourceM,
		sourceS: sourceS,
	}
}

func (r *ChatDbRepository) AddRoom(roomConfigInfo dm.ChatRoomConfigure) error {
	if roomConfigInfo.RoomOwnerUid == 0 || roomConfigInfo.RoomId == 0 || roomConfigInfo.RoomName == "" || roomConfigInfo.RedisSetCnt <= 0 {
		return errors.New("import param invalid")
	}

	_, err := r.sourceM.Insert(roomConfigInfo)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	return nil
}

func (r *ChatDbRepository) RemoveRoom(roomId int64) error {
	roomConfig := &dm.ChatRoomConfigure{
		Status: dm.RoomStatusDelete,
	}
	infectRow, err := r.sourceM.Where("room_id = ?", roomId).Cols("status").Update(roomConfig)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	if infectRow == 0 {
		return errors.New("invalid roomId")
	}
	return nil
}

func (r *ChatDbRepository) GetRoomConfigByRoomId(roomId int64) (info dm.ChatRoomConfigure, found bool) {
	found, err := r.sourceS.Where("room_id = ?", roomId).Get(&info)
	if err != nil {
		logger.Err(logType, err.Error())
		return dm.ChatRoomConfigure{}, false
	}
	return info, found
}

func (r *ChatDbRepository) GetRoomConfigByUid(uid int64) (info dm.ChatRoomConfigure, found bool) {
	found, err := r.sourceS.Where("uid = ?", uid).Get(&info)
	if err != nil {
		logger.Err(logType, err.Error())
		return dm.ChatRoomConfigure{}, false
	}
	return info, found
}

func (r *ChatDbRepository) ChangeRoomWorkStatus(uid int64, oriWorkStatus int8, newWorkStatus int8, needCheck bool) error {
	roomConfig := &dm.ChatRoomConfigure{
		WorkStatus: newWorkStatus,
	}
	var infectRow int64
	var err error
	if needCheck {
		infectRow, err = r.sourceM.Where("uid = ? and work_status = ?", uid, oriWorkStatus).Cols("work_status").Update(roomConfig)
	} else {
		infectRow, err = r.sourceM.Where("uid = ?", uid, oriWorkStatus).Cols("work_status").Update(roomConfig)
	}
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	if infectRow == 0 {
		return errors.New("invalid roomId")
	}
	return nil
}
