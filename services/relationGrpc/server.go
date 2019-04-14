package relationGrpc

import(
"context"
//"fmt"

dm "github.com/bluesky1024/goMblog/datamodels"
pb "github.com/bluesky1024/goMblog/services/relationGrpc/relationProto"
"github.com/bluesky1024/goMblog/tools/logger"
)

type relationService struct {
	client pb.RelationClient
}

func (s *relationService) GetFollowsByUid(uid int64, page int, pageSize int) (follows []dm.FollowInfo, cnt int64){
	req := pb.UidReq{
		Uid:uid,
		Page:int32(page),
		PageSize:int32(pageSize),
	}

	ctx := getSign(context.Background(),&req)
	res,err := s.client.GetFollowsByUid(ctx,&req)
	if err != nil {
		logger.Err(logType,err.Error())
		return follows,0
	}

	follows = make([]dm.FollowInfo,len(res.Follows))
	for ind,followInfo := range res.Follows {
		follows[ind] = dm.FollowInfo{
			Uid:followInfo.Uid,
			FollowUid:followInfo.FollowUid,
			Status:int8(followInfo.Status),
			GroupId:followInfo.GroupId,
		}
	}

	return follows,int64(len(res.Follows))
}

func (s *relationService) GetFansByUid(uid int64, page int, pageSize int) (fans []dm.FanInfo, cnt int64){
	req := pb.UidReq{
		Uid:uid,
		Page:int32(page),
		PageSize:int32(pageSize),
	}

	ctx := getSign(context.Background(),&req)
	res,err := s.client.GetFansByUid(ctx,&req)
	if err != nil {
		logger.Err(logType,err.Error())
		return fans,0
	}

	fans = make([]dm.FanInfo,len(res.Fans))
	for ind,fanInfo := range res.Fans {
		fans[ind] = dm.FanInfo{
			Uid:fanInfo.Uid,
			FanUid:fanInfo.FanUid,
			Status:int8(fanInfo.Status),
		}
	}

	return fans,int64(len(res.Fans))
}

func (s *relationService) GetGroupsByUid(uid int64)(groups []dm.FollowGroup,cnt int64){
	req := pb.UidReq{
		Uid:uid,
	}

	ctx := getSign(context.Background(),&req)
	res,err := s.client.GetGroupsByUid(ctx,&req)
	if err != nil {
		logger.Err(logType,err.Error())
		return groups,0
	}

	groups = make([]dm.FollowGroup,len(res.Groups))
	for ind,groupInfo := range res.Groups {
		groups[ind] = dm.FollowGroup{
			Uid:groupInfo.Uid,
			GroupName:groupInfo.GroupName,
			Status:int8(groupInfo.Status),
		}
	}

	return groups,int64(len(res.Groups))
}
