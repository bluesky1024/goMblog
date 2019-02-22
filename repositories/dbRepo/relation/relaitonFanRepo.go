package relationDbRepo

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func (r *RelationDbRepository) addFan(uid int64,uidFan int64,isFriend int8) bool {
	relation := dm.FanInfo{
		Uid:uid,
		FanUid:uidFan,
		Status:dm.FollowStatusNormal,
		IsFriend:isFriend,
	}
	_ ,err := r.sourceM.Table(getFanTableName(uid)).Insert(relation)
	if err != nil{
		logger.Err(logType,err.Error())
		return false
	}
	return true
}

func (r *RelationDbRepository) updateFan(info dm.FanInfo) bool {
	if info.Uid == 0 || info.FanUid == 0{
		return false
	}
	_ ,err := r.sourceM.Table(getFanTableName(info.Uid)).Cols("status","update_time","is_friend").Where("uid = ? and fan_uid = ?",info.Uid,info.FanUid).Update(&info)
	if err != nil {
		logger.Err(logType,err.Error())
		return false
	}
	return true
}

func (r *RelationDbRepository) SelectFanByUid(uid int64,uidFan int64) (relation dm.FanInfo, found bool){
	found,err := r.sourceS.Table(getFanTableName(uid)).Where("uid = ?",uid).And("fan_uid = ?",uidFan).Get(&relation)
	if err != nil {
		return dm.FanInfo{},false
	}
	return relation,found
}

func (r *RelationDbRepository) SelectMultiFansByUid(uid int64,page int,pageSize int) (infos []dm.FollowInfo,cnt int64){
	start := (page-1)*pageSize
	cnt, err := r.sourceS.Table(getFanTableName(uid)).Where("uid = ?",uid).Limit(pageSize,start).Asc("id").FindAndCount(&infos)
	if err!= nil {
		logger.Err(logType,err.Error())
		return infos,0
	}
	return infos,cnt
}

func (r *RelationDbRepository) AddOrUpdateFan(uid int64,uidFan int64) bool {
	//搜索是否存在关系
	relationA,foundA := r.SelectFanByUid(uid,uidFan)

	//若已存在,直接返回
	if foundA && (relationA.Status == dm.FollowStatusNormal) {
		return true
	}

	//数据更新
	var res bool
	if foundA {
		relationA.Status = dm.FollowStatusNormal
		res = r.updateFan(relationA)
	}else{
		res = r.addFan(uid,uidFan,dm.IsFriendFalse)
	}
	return res
}

func (r *RelationDbRepository) DeleteFan(uid int64,uidFan int64) bool {
	//搜索是否存在关系
	relationA,foundA := r.SelectFanByUid(uid,uidFan)

	//若不存在,直接返回
	if !foundA || (relationA.Status == dm.FollowStatusDelete) {
		return true
	}

	relationA.Status = dm.FollowStatusDelete
	return r.updateFan(relationA)
}