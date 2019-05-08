package mblogService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"strconv"
)

func TransMblogInfoToView(mblogs []dm.MblogInfo) (mblogsView []dm.MblogInfoView) {
	mblogsView = make([]dm.MblogInfoView, len(mblogs))
	for ind, mblog := range mblogs {
		mblogsView[ind] = dm.MblogInfoView{
			Mid:        strconv.FormatInt(mblog.Mid, 10),
			Uid:        strconv.FormatInt(mblog.Uid, 10),
			Content:    mblog.Content,
			OriginMid:  strconv.FormatInt(mblog.OriginMid, 10),
			OriginUid:  strconv.FormatInt(mblog.OriginUid, 10),
			TransCnt:   mblog.TransCnt,
			LikesCnt:   mblog.LikesCnt,
			CommentCnt: mblog.CommentCnt,
			Status:     mblog.Status,
			ReadAble:   mblog.ReadAble,
			CreateTime: mblog.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: mblog.UpdateTime.Format("2006-01-02 15:04:05"),
		}
	}
	return mblogsView
}
