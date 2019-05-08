package mblogGrpc

import (
	"context"
	dm "github.com/bluesky1024/goMblog/datamodels"
	pb "github.com/bluesky1024/goMblog/services/mblogGrpc/mblogProto"
	"github.com/bluesky1024/goMblog/tools/logger"
)

type mblogService struct {
	client pb.MblogClient
}

func (s *mblogService) Create(uid int64, content string, readAble int8, originUid int64, originMid int64) (mblog dm.MblogInfo, err error) {
	req := pb.MblogInfo{
		Uid:       uid,
		Content:   content,
		ReadAble:  int32(readAble),
		OriginUid: originUid,
		OriginMid: originMid,
	}

	ctx := getSign(context.Background(), &req)
	res, err := s.client.Create(ctx, &req)

	if err != nil {
		logger.Err(logType, err.Error())
		return mblog, err
	}

	mblog = dm.MblogInfo{
		Uid:       res.Uid,
		Content:   res.Content,
		ReadAble:  int8(res.ReadAble),
		OriginUid: res.OriginUid,
		OriginMid: res.OriginMid,
	}
	return mblog, nil
}

func (s *mblogService) GetMultiByMids(mids []int64) (mblogs map[int64]dm.MblogInfo) {
	req := pb.MidsReq{
		Mid: mids,
	}
	ctx := getSign(context.Background(), &req)
	res, err := s.client.GetMultiByMids(ctx, &req)
	if err != nil {
		logger.Err(logType, err.Error())
		return mblogs
	}
	mblogs = make(map[int64]dm.MblogInfo, len(res.MblogInfo))
	for _, v := range res.MblogInfo {
		mblogs[v.Mid] = dm.MblogInfo{
			Uid:       v.Uid,
			Mid:       v.Mid,
			Content:   v.Content,
			OriginUid: v.OriginUid,
			OriginMid: v.OriginMid,
		}
	}
	return mblogs
}

func (s *mblogService) GetNormalByUid(uid int64, page int, pageSize int, readAble []int8, startTime int64, endTime int64) (mblogs []dm.MblogInfo, cnt int64) {
	readAbleNew := make([]int32, len(readAble))
	for ind, v := range readAble {
		readAbleNew[ind] = int32(v)
	}
	req := pb.UidReq{
		Uid:       uid,
		Page:      int32(page),
		PageSize:  int32(pageSize),
		ReadAble:  readAbleNew,
		StartTime: startTime,
		EndTime:   endTime,
	}
	ctx := getSign(context.Background(), &req)
	res, err := s.client.GetNormalByUid(ctx, &req)

	if err != nil {
		logger.Err(logType, err.Error())
		return mblogs, 0
	}
	mblogs = make([]dm.MblogInfo, len(res.MblogInfo))
	for ind, v := range res.MblogInfo {
		mblogs[ind] = dm.MblogInfo{
			Uid:       v.Uid,
			Mid:       v.Mid,
			Content:   v.Content,
			OriginUid: v.OriginUid,
			OriginMid: v.OriginMid,
		}
	}
	return mblogs, res.Cnt
}
