package userService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
)

func (s *userService) HandleFollowMsg(msg dm.FollowMsg) (err error) {
	////更新
	//s.repo.AddOrUpdateFan(msg.FollowUid, msg.Uid)

	return err
}

func (s *userService) HandleUnFollowMsg(msg dm.FollowMsg) (err error) {
	////删除粉丝表记录
	//s.repo.DeleteFan(msg.FollowUid, msg.Uid)

	return err
}
