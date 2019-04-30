package userService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"strconv"
)

func TransUserInfoToView(users []dm.User) (usersView []dm.UserView) {
	usersView = make([]dm.UserView, len(users))
	for ind, user := range users {
		usersView[ind] = dm.UserView{
			Uid:          strconv.FormatInt(user.Uid, 10),
			NickName:     user.NickName,
			Telephone:    user.Telephone,
			Email:        user.Email,
			ProfileImage: user.ProfileImage,
			FollowsCount: strconv.FormatInt(user.FollowsCount, 10),
			FriendsCount: strconv.FormatInt(user.FriendsCount, 10),
			CreateTime:   user.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:   user.UpdateTime.Format("2006-01-02 15:04:05"),
		}
	}
	return usersView
}
