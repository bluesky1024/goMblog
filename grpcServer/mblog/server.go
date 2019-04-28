package main

import (
	"context"
	"errors"
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

func (s *MblogService) GetByUid(ctx context.Context, uidReq *pb.UidReq) (res *pb.MultiMblogs, err error) {
	mblogs := s.serv.GetByUid(uidReq.Uid, int(uidReq.Page), int(uidReq.PageSize))

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

	return res, err
}

func (s *MblogService) GetByMid(ctx context.Context, midReq *pb.MidReq) (res *pb.MblogInfo, err error) {
	mblog, found := s.serv.GetByMid(midReq.Mid)
	if !found {
		return res, errors.New("not found mblog")
	}

	return &pb.MblogInfo{
		Mid:       mblog.Mid,
		Uid:       mblog.Uid,
		Content:   mblog.Content,
		OriginMid: mblog.OriginMid,
		OriginUid: mblog.OriginUid,
	}, nil
}
