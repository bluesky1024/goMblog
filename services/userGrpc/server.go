package userGrpc

import(
	"context"
	"fmt"

	dm "github.com/bluesky1024/goMblog/datamodels"
	pb "github.com/bluesky1024/goMblog/services/userGrpc/userProto"
	"github.com/bluesky1024/goMblog/tools/logger"
)

type userService struct {
	client pb.MblogUserClient
}

func (s *userService) Create(nickname string, password string, telephone string, email string) (user dm.User, err error) {
	req := pb.User{
		NickName:nickname,
		Password:password,
		Telephone:telephone,
		Email:email,
	}

	ctx := getSign(context.Background(),&req)
	res,err := s.client.Create(ctx,&req)

	if err != nil {
		logger.Err(logType,err.Error())
		return user,err
	}

	fmt.Println(res.Uid)

	user = dm.User{
		Id:res.Id,
		Uid:res.Uid,
		NickName:res.NickName,
		Password:res.Password,
		Telephone:res.Telephone,
		Email:res.Email,
		ProfileImage:res.ProfileImages,
		FollowsCount:res.FollowsCount,
		FriendsCount:res.FriendsCount,
		//CreateTime:time.Unix(res.CreateTime.Seconds,0),
		//UpdateTime:time.Unix(res.UpdateTime.Seconds,0),
	}
	return user,nil
}