package feedService

func (s *feedService) HandleFollowMsg(msg dm.FollowMsg) (err error) {

	return err
}

func (s *feedService) HandleUnFollowMsg(msg dm.FollowMsg) (err error) {
	//删除粉丝表记录
	s.repo.DeleteFan(msg.FollowUid, msg.Uid)

	return err
}
