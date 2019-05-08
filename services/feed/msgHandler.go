package feedService

import (
	"errors"
	"fmt"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"time"
)

func (f *feedService) HandleFollowMsg(msg dm.FollowMsg) (err error) {
	//检查两者关系
	//relation := f.relationSrv.check_relation(msg.Uid, msg.FollowUid)
	readAble := []int8{dm.MblogReadAblePublic}

	need, err := f.checkNeedPullUserMblog(msg.FollowUid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	if !need {
		return nil
	}

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

func (f *feedService) HandleSetGroupMsg(msg dm.SetGroupMsg) (err error) {
	need, err := f.checkNeedPullUserMblog(msg.FollowUid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	if !need {
		return nil
	}

	//由于一个人只能属于一个关注分组,所以需要先把之前的分组feed内的个人微博数据移除，再给新的分组内的feed增加数据
	//移除原始分组feed数据
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
			f.appendNewMid(msg.Uid, msg.GroupId, mblog.Mid)
			mids[ind] = mblog.Mid
		}
		if msg.OriGroupId != 0 {
			f.feedRdRepo.RemoveMids(msg.Uid, msg.OriGroupId, mids)
		}

		page++
	}

	return nil
}

func (f *feedService) HandleMblogNewMsg(msg dm.MblogNewMsg) (err error) {
	//判断微博的可见属性,如果是私人可见的，不进行feed处理(目前暂未考虑朋友可见这一属性，后期待完善)
	if msg.ReadAble == dm.MblogReadAblePerson || msg.ReadAble == dm.MblogReadAbleFriend {
		return nil
	}

	need, err := f.checkNeedPullUserMblog(msg.Uid)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	if !need {
		return nil
	}

	//提取该用户的粉丝列表,给所有粉丝的feed池新增数据
	page := 1
	pageSize := 50
	for {
		fans, _ := f.relationSrv.GetFansByUid(msg.Uid, page, pageSize)
		if len(fans) == 0 {
			break
		}

		for _, fan := range fans {
			res := f.appendNewMid(fan.FanUid, 0, msg.Mid)
			fmt.Println("append res", res)
			if fan.GroupId != 0 {
				f.appendNewMid(fan.FanUid, fan.GroupId, msg.Mid)
			}
		}
		page++
	}

	return nil
}
