package relationService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func (s *relationService) GetFollowCntByUids(uids []int64) map[int64]int64 {
	res, err := s.rdRepo.GetFollowCnt(uids)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return res
}

func (s *relationService) GetFollowsByUid(uid int64, page int, pageSize int) (follows []dm.FollowInfo, totalCnt int64) {
	follows, _ = s.repo.SelectMultiFollowsByUid(uid, page, pageSize)
	cntMap, err := s.rdRepo.GetFollowCnt([]int64{uid})
	if len(cntMap) == 0 || err != nil {
		return follows, 0
	}
	return follows, cntMap[uid]
}

func (s *relationService) Follow(uid int64, uidFollow int64) bool {
	//修改follow表
	succ := s.repo.AddOrUpdateFollow(uid, uidFollow)

	//其他操作加入消息队列进行操作
	if succ {
		msg := dm.FollowMsg{
			Uid:       uid,
			FollowUid: uidFollow,
			GroupId:   0,
			Status:    dm.FollowStatusNormal,
		}
		s.sendFollowMsg(msg)
	}
	return succ
}

func (s *relationService) UnFollow(uid int64, uidFollow int64) bool {
	//获取原有follow属性
	followInfo, found := s.repo.SelectFollowByUid(uid, uidFollow)
	if !found {
		return false
	}
	//修改follow表
	succ := s.repo.DeleteFollow(uid, uidFollow)

	//其他操作加入消息队列进行操作
	if succ {
		msg := dm.FollowMsg{
			Uid:       uid,
			FollowUid: uidFollow,
			GroupId:   followInfo.GroupId,
			Status:    dm.FollowStatusDelete,
		}
		s.sendUnFollowMsg(msg)
	}
	return succ
}

func (s *relationService) CheckFollow(uidA int64, uidB int64) int {
	if uidA == 0 || uidB == 0 || uidA == uidB {
		return 0
	}

	info, found := s.repo.SelectFollowByUid(uidA, uidB)
	if !found || info.Status == dm.FollowStatusDelete {
		return 0
	}

	return 1
}

func (s *relationService) CheckRelation(uidA int64, uidB int64) int8 {
	if uidA == 0 || uidB == 0 {
		return dm.RelationNone
	}
	if uidA == uidB {
		return dm.RelationSelf
	}
	return dm.RelationNone
}
