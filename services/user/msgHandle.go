package userService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
)

//根据关注取关动作消息维护粉丝数的变更
func (s *userService) HandleFollowMsg(msg dm.FollowMsg) (err error) {
	//被关注用户粉丝数加一

	return err
}

func (s *userService) HandleUnFollowMsg(msg dm.FollowMsg) (err error) {
	//被取关用户粉丝数减一

	return err
}
