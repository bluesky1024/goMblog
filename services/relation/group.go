package relationService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
)

func (s *relationService) GetGroupsByUid(uid int64) (groups []dm.FollowGroup, cnt int64) {
	return s.repo.SelectMultiGroupsByUid(uid)
}

func (s *relationService) AddGroup(uid int64, groupName string) bool {
	return s.repo.AddOrUpdateGroup(uid, groupName)
}

func (s *relationService) DelGroup(uid int64, groupId int64) bool {
	return s.repo.DeleteGroupByUidAndGroupId(uid, groupId)
}

func (s *relationService) UpdateGroup(group dm.FollowGroup) bool {
	return s.repo.UpdateGroupById(group)
}

func (s *relationService) SetFollowGroup(uid int64, uidFollow int64, groupId int64) bool {
	res, _ := s.repo.UpdateFollowGroupByUid(uid, uidFollow, groupId)
	return res
}
