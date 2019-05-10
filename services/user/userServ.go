package userService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	ds "github.com/bluesky1024/goMblog/datasource/dbSource"
	"github.com/bluesky1024/goMblog/repositories/dbRepo/user"
	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/tools/logger"
)

var logType = "userService"

type UserServicer interface {
	Create(nickname string, password string, telephone string, email string) (user dm.User, err error)
	GetByNicknameAndPassword(nickname string, password string) (user dm.User, found bool)
	GetByUid(uid int64) (user dm.User, found bool)
	GetMultiByUids(uids []int64) (users map[int64]dm.User, err error)

	//	GetAll() []datamodels.User
	//	GetByID(id int64) (datamodels.User, bool)
	//	GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool)
	//	DeleteByID(id int64) bool

	//	Update(id int64, user datamodels.User) (datamodels.User, error)
	//	UpdatePassword(id int64, newPassword string) (datamodels.User, error)
	//	UpdateUsername(id int64, newUsername string) (datamodels.User, error)

	HandleFollowMsg(msg dm.FollowMsg) error
	HandleUnFollowMsg(msg dm.FollowMsg) error
}

type userService struct {
	repo *userDbRepo.UserDbRepository
}

// NewUserService returns the default user service.
func NewUserServicer() (s UserServicer, err error) {
	//id生成池初始化
	err = idGen.InitUidPool(3)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}

	//user服务仓库初始化
	userSourceM, err := ds.LoadUsers(true)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	userSourceR, err := ds.LoadUsers(false)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	userRepo := userDbRepo.NewUserRepository(userSourceM, userSourceR)

	return &userService{
		repo: userRepo,
	}, nil
}
