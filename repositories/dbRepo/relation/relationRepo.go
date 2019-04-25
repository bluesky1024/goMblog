package relationDbRepo

import (
	"errors"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"github.com/go-xorm/xorm"
)

var logType string = "relationDbRepo"

type RelationDbRepository struct {
	sourceM *xorm.Engine
	sourceS *xorm.Engine
}

func NewRelationRepository(sourceM *xorm.Engine, sourceS *xorm.Engine) *RelationDbRepository {
	return &RelationDbRepository{
		sourceM: sourceM,
		sourceS: sourceS,
	}
}

func (r *RelationDbRepository) addFollow(uid int64, uidFollow int64, isFriend int8) bool {
	relation := dm.FollowInfo{
		Uid:       uid,
		FollowUid: uidFollow,
		Status:    dm.FollowStatusNormal,
		IsFriend:  isFriend,
	}
	_, err := r.sourceM.Table(getFollowTableName(uid)).Insert(relation)
	if err != nil {
		logger.Err(logType, err.Error())
		return false
	}
	return true
}

func (r *RelationDbRepository) updateFollowByUid(info dm.FollowInfo) bool {
	if info.Uid == 0 || info.FollowUid == 0 {
		return false
	}
	_, err := r.sourceM.Table(getFollowTableName(info.Uid)).Cols("status", "update_time", "is_friend", "group_id").Where("uid = ? and follow_uid = ?", info.Uid, info.FollowUid).Update(&info)
	if err != nil {
		logger.Err(logType, err.Error())
		return false
	}
	return true
}

func (r *RelationDbRepository) SelectFollowByUid(uid int64, uidFollow int64) (relation dm.FollowInfo, found bool) {
	found, err := r.sourceS.Table(getFollowTableName(uid)).Where("uid = ?", uid).And("follow_uid = ?", uidFollow).Get(&relation)
	if err != nil {
		return dm.FollowInfo{}, false
	}
	return relation, found
}

func (r *RelationDbRepository) AddOrUpdateFollow(uid int64, uidFollow int64) bool {
	//搜索是否存在关系
	relationA, foundA := r.SelectFollowByUid(uid, uidFollow)

	//若已存在,直接返回
	if foundA && (relationA.Status == dm.FollowStatusNormal) {
		return true
	}

	//是否互关，该参数暂时不处理

	//数据更新
	var res bool
	if foundA {
		relationA = dm.FollowInfo{
			Uid:       uid,
			FollowUid: uidFollow,
			Status:    dm.FollowStatusNormal,
			IsFriend:  dm.IsFriendFalse,
			GroupId:   0,
		}
		res = r.updateFollowByUid(relationA)
	} else {
		res = r.addFollow(uid, uidFollow, dm.IsFriendFalse)
	}
	return res
}

func (r *RelationDbRepository) DeleteFollow(uid int64, uidFollow int64) bool {
	//搜索是否存在关系
	relationA, foundA := r.SelectFollowByUid(uid, uidFollow)

	//若不存在,直接返回
	if !foundA || (relationA.Status == dm.FollowStatusDelete) {
		return true
	}

	relationA = dm.FollowInfo{
		Uid:       uid,
		FollowUid: uidFollow,
		Status:    dm.FollowStatusDelete,
		IsFriend:  dm.IsFriendFalse,
		GroupId:   0,
	}
	return r.updateFollowByUid(relationA)
}

func (r *RelationDbRepository) UpdateFollowGroupByUid(uid int64, uidFollow int64, groupId int64) (res bool, err error) {
	//搜索是否存在关系
	relationA, foundA := r.SelectFollowByUid(uid, uidFollow)

	//若不存在,直接返回
	if !foundA || (relationA.Status == dm.FollowStatusDelete) {
		return false, errors.New("not found relation")
	}

	relationA.GroupId = groupId
	res = r.updateFollowByUid(relationA)
	return res, nil
}

func (r *RelationDbRepository) SelectMultiFollowsByUid(uid int64, page int, pageSize int) (infos []dm.FollowInfo, cnt int64) {
	start := (page - 1) * pageSize
	cnt, err := r.sourceS.Table(getFollowTableName(uid)).Where("uid = ?", uid).Limit(pageSize, start).Asc("id").FindAndCount(&infos)
	if err != nil {
		logger.Err(logType, err.Error())
		return infos, 0
	}
	return infos, cnt
}
