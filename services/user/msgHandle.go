package userService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

//根据关注取关动作消息维护粉丝数的变更
func (s *userService) HandleFollowMsg(msg dm.FollowMsg) (err error) {
	//被关注用户粉丝数加一
	_, err = s.redisRepo.AddFan(msg.FollowUid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	//关注用户关注数加一
	_, err = s.redisRepo.AddFollow(msg.Uid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	return nil
}

func (s *userService) HandleUnFollowMsg(msg dm.FollowMsg) (err error) {
	//被关注用户粉丝数减一
	_, err = s.redisRepo.LoseFan(msg.FollowUid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	//关注用户关注数减一
	_, err = s.redisRepo.LoseFollow(msg.Uid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	return nil
}
