package relationDbRepo

import "github.com/bluesky1024/goMblog/tools/logger"

func (r *RelationDbRepository) UpdateFanCnt(uid int64, uidFan int64, isFriend int8) bool {
	relation := dm.FanInfo{
		Uid:      uid,
		FanUid:   uidFan,
		Status:   dm.FollowStatusNormal,
		IsFriend: isFriend,
	}
	_, err := r.sourceM.Table(getFanTableName(uid)).Insert(relation)
	if err != nil {
		logger.Err(logType, err.Error())
		return false
	}
	return true
}

func (r *RelationDbRepository) UpdateFollowCnt(uid int64, uidFan int64, isFriend int8) bool {
	relation := dm.FanInfo{
		Uid:      uid,
		FanUid:   uidFan,
		Status:   dm.FollowStatusNormal,
		IsFriend: isFriend,
	}
	_, err := r.sourceM.Table(getFanTableName(uid)).Incr(relation)
	if err != nil {
		logger.Err(logType, err.Error())
		return false
	}
	return true
}
