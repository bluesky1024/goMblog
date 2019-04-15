package main

import (
	"context"
	pb "github.com/bluesky1024/goMblog/grpcServer/relation/relationProto"
	"github.com/bluesky1024/goMblog/services/relation"
)

func NewRelationService() *RelationService {
	return &RelationService{
		serv: relationService.NewRelationServicer(),
	}
}

type RelationService struct {
	serv relationService.RelationServicer
}

func (s *RelationService) GetFollowsByUid(ctx context.Context, req *pb.UidReq) (res *pb.MultiFollows, err error) {
	follows, _ := s.serv.GetFollowsByUid(req.Uid, int(req.Page), int(req.PageSize))

	//follows组包
	res = &pb.MultiFollows{}
	res.Follows = make([]*pb.FollowInfo, len(follows))
	for ind, follow := range follows {
		res.Follows[ind] = &pb.FollowInfo{
			Uid:       follow.Uid,
			FollowUid: follow.FollowUid,
			Status:    int32(follow.Status),
			GroupId:   follow.GroupId,
		}
	}
	return res, err
}

func (s *RelationService) GetFansByUid(ctx context.Context, req *pb.UidReq) (res *pb.MultiFans, err error) {
	fans, _ := s.serv.GetFansByUid(req.Uid, int(req.Page), int(req.PageSize))

	//fans组包
	res = &pb.MultiFans{}
	res.Fans = make([]*pb.FanInfo, len(fans))
	for ind, fan := range fans {
		res.Fans[ind] = &pb.FanInfo{
			Uid:    fan.Uid,
			FanUid: fan.FanUid,
			Status: int32(fan.Status),
		}
	}
	return res, err
}

func (s *RelationService) GetGroupsByUid(ctx context.Context, req *pb.UidReq) (res *pb.MultiFollowGroups, err error) {
	groups, _ := s.serv.GetGroupsByUid(req.Uid)

	//groups组包
	res = &pb.MultiFollowGroups{}
	res.Groups = make([]*pb.FollowGroup, len(groups))
	for ind, group := range groups {
		res.Groups[ind] = &pb.FollowGroup{
			Id:        group.Id,
			Uid:       group.Uid,
			GroupName: group.GroupName,
			Status:    int32(group.Status),
		}
	}
	return res, err
}
