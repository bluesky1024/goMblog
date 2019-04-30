package feedService

import (
	"errors"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"time"
)

func (f *feedService) HandleFollowMsg(msg dm.FollowMsg) (err error) {
	//检查两者关系
	//relation := f.relationSrv.check_relation(msg.Uid, msg.FollowUid)
	readAble := []int8{dm.MblogReadAblePublic}

	//循环获取关注人近30天的微博mid
	endTime := time.Now().Unix()
	startTime := endTime - 30*60*60*24
	page := 1
	pageSize := 50
	for {
		mblogs, _ := f.mblogSrv.GetNormalByUid(msg.FollowUid, page, pageSize, readAble, startTime, 0)
		if len(mblogs) == 0 {
			break
		}

		for _, mblog := range mblogs {
			f.appendNewMid(msg.Uid, 0, mblog.Mid)
		}
		page++
	}
	return nil
}

func (f *feedService) HandleUnFollowMsg(msg dm.FollowMsg) (err error) {
	//检查两者关系
	//relation := f.relationSrv.check_relation(msg.Uid, msg.FollowUid)
	readAble := []int8{dm.MblogReadAblePublic}

	//循环获取关注人近30天的微博mid
	endTime := time.Now().Unix()
	startTime := endTime - 30*60*60*24
	page := 1
	pageSize := 50
	for {
		mblogs, _ := f.mblogSrv.GetNormalByUid(msg.FollowUid, page, pageSize, readAble, startTime, 0)
		if len(mblogs) == 0 {
			break
		}

		mids := make([]int64, len(mblogs))
		for ind, mblog := range mblogs {
			mids[ind] = mblog.Mid
		}
		f.feedRdRepo.RemoveMids(msg.Uid, 0, mids)
		f.feedRdRepo.RemoveMids(msg.Uid, msg.GroupId, mids)
		page++
	}
	return nil
}

func (f *feedService) HandleMblogNewMsg(msg dm.MblogNewMsg) (err error) {
	//判断该用户的属性(若粉丝数太多，feed采用推模式，不用给每个粉丝的feed池里逐一添加)
	userInfo, err := f.userSrv.GetByUid(msg.Uid)
	if err != nil {
		err = errors.New("handle mblog new msg, uid is invalid")
		logger.Err(logType, err.Error())
		return err
	}

	//判断标准，粉丝数大于xx
	if userInfo.FriendsCount >= int64(dm.FansCntGate) {
		return nil
	}

	//提取该用户的粉丝列表,给所有粉丝的feed池新增数据
	page := 1
	pageSize := 50
	for {

	}

	return nil
}
