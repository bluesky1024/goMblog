package userGrpc

import (
	"context"
	//"fmt"

	dm "github.com/bluesky1024/goMblog/datamodels"
	pb "github.com/bluesky1024/goMblog/services/userGrpc/userProto"
	"github.com/bluesky1024/goMblog/tools/logger"
)

type userService struct {
	client pb.MblogUserClient
}

func (s *userService) Create(nickname string, password string, telephone string, email string) (user dm.User, err error) {
	req := pb.User{
		NickName:  nickname,
		Password:  password,
		Telephone: telephone,
		Email:     email,
	}

	ctx := getSign(context.Background(), &req)
	res, err := s.client.Create(ctx, &req)

	if err != nil {
		logger.Err(logType, err.Error())
		return user, err
	}

	user = dm.User{
		//Id:res.Id,
		Uid:      res.Uid,
		NickName: res.NickName,
		//Password:res.Password,
		Telephone: res.Telephone,
		Email:     res.Email,
		//ProfileImage:res.ProfileImages,
		//FollowsCount:res.FollowsCount,
		//FriendsCount:res.FriendsCount,
		//CreateTime:time.Unix(res.CreateTime.Seconds,0),
		//UpdateTime:time.Unix(res.UpdateTime.Seconds,0),
	}
	return user, nil
}

func (s *userService) GetByUid(uid int64) (user dm.User, err error) {
	req := pb.Uid{
		Uid: uid,
	}

	ctx := getSign(context.Background(), &req)
	res, err := s.client.GetByUid(ctx, &req)

	if err != nil {
		logger.Err(logType, err.Error())
		return user, err
	}

	user = dm.User{
		Uid:       res.Uid,
		NickName:  res.NickName,
		Telephone: res.Telephone,
		Email:     res.Email,
	}
	return user, nil
}

func (s *userService) GetByNicknameAndPassword(nickName string, passWord string) (user dm.User, err error) {
	req := pb.User{
		NickName: nickName,
		Password: passWord,
	}

	ctx := getSign(context.Background(), &req)
	res, err := s.client.GetByNicknameAndPassword(ctx, &req)

	if err != nil {
		logger.Err(logType, err.Error())
		return user, err
	}

	user = dm.User{
		Uid:       res.Uid,
		NickName:  res.NickName,
		Telephone: res.Telephone,
		Email:     res.Email,
	}
	return user, nil
}

func (s *userService) GetMultiByUids(uids []int64) (users []dm.User, err error) {
	req := pb.Uids{}
	req.SingleUid = make([]*pb.Uid, len(uids))
	for ind, uid := range uids {
		req.SingleUid[ind] = &pb.Uid{Uid: uid}
	}

	ctx := getSign(context.Background(), &req)
	res, err := s.client.GetMultiByUids(ctx, &req)

	if err != nil {
		logger.Err(logType, err.Error())
		return users, err
	}

	users = make([]dm.User, len(res.UserInfo))
	for ind, userInfo := range res.UserInfo {
		users[ind] = dm.User{
			Uid:       userInfo.Uid,
			NickName:  userInfo.NickName,
			Telephone: userInfo.Telephone,
			Email:     userInfo.Email,
		}
	}
	return users, nil
}
