package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/bluesky1024/goMblog/grpcServer/user/userProto"
	userServ "github.com/bluesky1024/goMblog/services/user"
	"github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func NewMblogUserService() *MblogUserService {
	err := idGenerate.InitUidPool(3)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil
	}
	s, err := userServ.NewUserServicer()
	if err != nil {
		logger.Err(logType, err.Error())
		return nil
	}

	return &MblogUserService{
		serv: s,
	}
}

type MblogUserService struct {
	serv userServ.UserServicer
}

func (s *MblogUserService) Create(ctx context.Context, user *pb.User) (res *pb.User, err error) {
	newUser, err := s.serv.Create(user.NickName, user.Password, user.Telephone, user.Email)
	if err != nil {
		return res, err
	}

	return &pb.User{
		Uid:       newUser.Uid,
		NickName:  newUser.NickName,
		Telephone: newUser.Telephone,
		Email:     newUser.Email,
	}, nil
}

func (s *MblogUserService) GetByUid(ctx context.Context, uid *pb.Uid) (res *pb.User, err error) {
	userInfo, found := s.serv.GetByUid(uid.Uid)
	if !found {
		return res, errors.New("not found")
	}
	return &pb.User{
		Uid:       userInfo.Uid,
		NickName:  userInfo.NickName,
		Telephone: userInfo.Telephone,
		Email:     userInfo.Email,
	}, nil
}

func (s *MblogUserService) GetByNicknameAndPassword(ctx context.Context, user *pb.User) (res *pb.User, err error) {
	userInfo, found := s.serv.GetByNicknameAndPassword(user.NickName, user.Password)
	if !found {
		return res, errors.New("not found")
	}
	return &pb.User{
		Uid:       userInfo.Uid,
		NickName:  userInfo.NickName,
		Telephone: userInfo.Telephone,
		Email:     userInfo.Email,
	}, nil
}

func (s *MblogUserService) GetMultiByUids(ctx context.Context, uids *pb.Uids) (res *pb.MultiUsers, err error) {
	//uids拆包
	uidSlice := make([]int64, len(uids.SingleUid))
	for ind, uid := range uids.SingleUid {
		uidSlice[ind] = uid.Uid
	}
	users, err := s.serv.GetMultiByUids(uidSlice)
	if err != nil {
		return res, err
	}

	//users组包
	res = &pb.MultiUsers{}
	fmt.Println(len(users))
	res.UserInfo = make([]*pb.User, len(users))
	cnt := 0
	for _, user := range users {
		res.UserInfo[cnt] = &pb.User{
			Uid:       user.Uid,
			NickName:  user.NickName,
			Telephone: user.Telephone,
			Email:     user.Email,
		}
		cnt++
	}

	return res, err
}
