package chatRdRepo

import (
	"encoding/json"
	"github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func (r *ChatRbRepository) SetRoomConfig(uid int64, roomConfig datamodels.ChatRoomConfigure) error {
	key, expireTimeDuration := getRoomConfigSetInfo(uid)
	configInfo, _ := json.Marshal(&roomConfig)
	_, err := r.redisPool.Set(key, configInfo, expireTimeDuration).Result()
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}

func (r *ChatRbRepository) GetRoomConfig(uid int64) (datamodels.ChatRoomConfigure, error) {
	key, _ := getRoomConfigSetInfo(uid)
	res, err := r.redisPool.Get(key).Result()
	config := new(datamodels.ChatRoomConfigure)
	if err != nil {
		logger.Err(logType, err.Error())
		return *config, err
	}
	err = json.Unmarshal([]byte(res), config)
	if err != nil {
		logger.Err(logType, err.Error())
		return *config, err
	}
	return *config, nil
}

func (r *ChatRbRepository) DelRoomConfig(uid int64) (succ bool, err error) {
	key, _ := getRoomConfigSetInfo(uid)
	res, err := r.redisPool.Del(key).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return false, err
	}
	if res != 0 {
		return true, nil
	}
	return false, nil
}
