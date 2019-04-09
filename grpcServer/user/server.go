package main

import (
	"context"
	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"
	//"fmt"
	userServ "github.com/bluesky1024/goMblog/services/user"
	pb "github.com/bluesky1024/goMblog/grpcServer/user/userProto"
)

func NewMblogUserService() *MblogUserService {
	idGen.InitUidPool(3)
	return &MblogUserService{
		serv:userServ.NewUserServicer(),
	}
}

type MblogUserService struct{
	serv userServ.UserServicer
}

func (s *MblogUserService) Create(ctx context.Context, user *pb.User) (res *pb.User, err error) {
	newUser,err := s.serv.Create(user.NickName,user.Password,user.Telephone,user.Email)
	if err != nil {
		return res,err
	}

	return &pb.User{
		Uid:newUser.Uid,
		NickName:newUser.NickName,
		Telephone:newUser.Telephone,
		Email:newUser.Email,
	},nil
}

func (s *MblogUserService) GetByUid(ctx context.Context, uid *pb.Uid) (res *pb.User, err error) {
	return res, err
}

func (s *MblogUserService) GetByNicknameAndPassword(ctx context.Context, user *pb.User) (res *pb.User, err error) {
	return res, err
}

func (s *MblogUserService) GetMultiByUids(ctx context.Context, uids *pb.Uids) (res *pb.MultiUsers, err error) {
	return res, err
}
