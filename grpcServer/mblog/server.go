package main

import (
	"context"
	pb "github.com/bluesky1024/goMblog/grpcServer/mblog/mblogProto"
	"github.com/bluesky1024/goMblog/services/mblog"
	"github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func NewMblogService() *MblogService {
	err := idGenerate.InitMidPool(10)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil
	}

	s, err := mblogService.NewMblogServicer()
	if err != nil {
		logger.Err(logType, err.Error())
		return nil
	}

	return &MblogService{
		serv: s,
	}
}

type MblogService struct {
	serv mblogService.MblogServicer
}

func (s *MblogService) Create(ctx context.Context, mblogInfo *pb.MblogInfo) (res *pb.MblogInfo, err error) {
	res = new(pb.MblogInfo)
	mblog, err := s.serv.Create(mblogInfo.Uid, mblogInfo.Content, int8(mblogInfo.ReadAble), mblogInfo.OriginUid, mblogInfo.OriginMid)
	if err != nil {
		return res, err
	}

	return &pb.MblogInfo{
		Mid:       mblog.Mid,
		Uid:       mblog.Uid,
		Content:   mblog.Content,
		OriginMid: mblog.OriginMid,
		OriginUid: mblog.OriginUid,
	}, nil
}

func (s *MblogService) GetNormalByUid(ctx context.Context, uidReq *pb.UidReq) (res *pb.MultiMblogs, err error) {
	res = new(pb.MultiMblogs)
	readAble := make([]int8, len(uidReq.ReadAble))
	for ind, v := range uidReq.ReadAble {
		readAble[ind] = int8(v)
	}
	mblogs, cnt := s.serv.GetNormalByUid(uidReq.Uid, int(uidReq.Page), int(uidReq.PageSize), readAble, uidReq.StartTime, uidReq.EndTime)

	res.MblogInfo = make([]*pb.MblogInfo, len(mblogs))

	for ind, mblog := range mblogs {
		res.MblogInfo[ind] = &pb.MblogInfo{
			Mid:       mblog.Mid,
			Uid:       mblog.Uid,
			Content:   mblog.Content,
			OriginMid: mblog.OriginMid,
			OriginUid: mblog.OriginUid,
		}
	}
	res.Cnt = cnt

	return res, err
}

func (s *MblogService) GetMultiByMids(ctx context.Context, midsReq *pb.MidsReq) (res *pb.MultiMblogs, err error) {
	res = new(pb.MultiMblogs)
	mblogMap := s.serv.GetMultiByMids(midsReq.Mid)

	res.MblogInfo = make([]*pb.MblogInfo, len(mblogMap))

	ind := 0
	for _, mblog := range mblogMap {
		res.MblogInfo[ind] = &pb.MblogInfo{
			Mid:       mblog.Mid,
			Uid:       mblog.Uid,
			Content:   mblog.Content,
			OriginMid: mblog.OriginMid,
			OriginUid: mblog.OriginUid,
		}
		ind++
	}
	res.Cnt = int64(ind)

	return res, err
}
