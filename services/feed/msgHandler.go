package feedService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
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
