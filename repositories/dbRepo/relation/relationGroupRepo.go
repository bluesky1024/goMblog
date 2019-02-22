package relationDbRepo

import (
	"fmt"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func (r *RelationDbRepository) addGroup(uid int64,groupName string) bool {
	if uid == 0 || groupName == "" {
		return false
	}

	group := dm.FollowGroup{
		Uid:uid,
		GroupName:groupName,
		Status:dm.GroupStatusNormal,
	}
	_ ,err := r.sourceM.Table(getFollowGroupTableName(uid)).Insert(group)
	if err != nil{
		logger.Err(logType,err.Error())
		return false
	}
	return true
}

func (r *RelationDbRepository) UpdateGroupById(info dm.FollowGroup) bool {
	if info.Id == 0 || info.Uid == 0 || info.GroupName == ""{
		return false
	}
	_ ,err := r.sourceM.Table(getFollowGroupTableName(info.Uid)).Cols("group_name","status").Where("id = ? and uid = ?",info.Id,info.Uid).Update(&info)
	if err != nil {
		logger.Err(logType,err.Error())
		return false
	}
	return true
}

func (r *RelationDbRepository) SelectGroupByUidAndGroupId(uid int64,groupId int64) (info dm.FollowGroup,found bool) {
	found,err := r.sourceS.Table(getFollowGroupTableName(uid)).Where("id = ? and uid = ?",groupId,uid).Get(&info)
	if err != nil {
		logger.Err(logType,err.Error())
		return dm.FollowGroup{},false
	}
	return info,found
}

func (r *RelationDbRepository) SelectGroupByUidAndName(uid int64,groupName string) (info dm.FollowGroup,found bool) {
	found,err := r.sourceS.Table(getFollowGroupTableName(uid)).Where("uid = ? and group_name = ?",uid,groupName).Get(&info)
	if err != nil {
		logger.Err(logType,err.Error())
		return dm.FollowGroup{},false
	}
	return info,found
}

func (r *RelationDbRepository) SelectMultiGroupsByUid(uid int64) (infos []dm.FollowGroup,cnt int64) {
	cnt, err := r.sourceS.Table(getFollowGroupTableName(uid)).Where("uid = ? and status = ?",uid,dm.GroupStatusNormal).Asc("id").FindAndCount(&infos)
	if err!= nil {
		logger.Err(logType,err.Error())
		return infos,0
	}
	return infos,cnt
}

func (r *RelationDbRepository) AddOrUpdateGroup(uid int64,groupName string) bool {
	//搜索是否存在
	group,found := r.SelectGroupByUidAndName(uid,groupName)

	//若已存在,直接返回
	if found && (group.Status == dm.GroupStatusNormal) {
		return true
	}

	//数据更新
	var res bool
	if found {
		group.Status = dm.GroupStatusNormal
		res = r.UpdateGroupById(group)
	}else{
		res = r.addGroup(uid,groupName)
	}
	return res
}

func (r *RelationDbRepository) DeleteGroupByUidAndGroupId(uid int64,groupId int64) bool {
	if uid == 0 || groupId == 0 {
		return false
	}
	//搜索是否存在关系
	group,found := r.SelectGroupByUidAndGroupId(uid,groupId)

	fmt.Println("sel",group)

	//若不存在,直接返回
	if !found || (group.Status == dm.GroupStatusDelete) {
		return true
	}

	group.Status = dm.GroupStatusDelete
	fmt.Println("del",group)
	return r.UpdateGroupById(group)
}