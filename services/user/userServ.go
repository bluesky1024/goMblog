package userService

import (
	"errors"
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

// Create inserts a new User,
// the userPassword is the client-typed password
// it will be hashed before the insertion to our repository.
func (s *userService) Create(nickname string, password string, telephone string, email string) (dm.User, error) {
	if nickname == "" || password == "" || telephone == "" || email == "" {
		return dm.User{}, errors.New("unable to create this user")
	}
	//生成password对应hash值
	hashed, err := dm.GeneratePassword(password)
	if err != nil {
		return dm.User{}, err
	}

	insertUser := dm.User{
		NickName:  nickname,
		Password:  hashed,
		Telephone: telephone,
		Email:     email,
	}

	//生成专属uid
	insertUser.Uid, err = idGen.GenUidId()
	if err != nil {
		return dm.User{}, err
	}

	affect, err := s.repo.Insert(insertUser)
	if err != nil || affect == 0 {
		return dm.User{}, err
	}
	return insertUser, nil
}

// GetByID returns a user based on its id.
func (s *userService) GetByUid(uid int64) (user dm.User, found bool) {
	if uid < 0 {
		return dm.User{}, false
	}

	user, found = s.repo.SelectByUid(uid)

	//验证密码是否正确
	if found {
		return user, true
	}
	return dm.User{}, false
}

func (s *userService) GetMultiByUids(uids []int64) (users map[int64]dm.User, err error) {
	return s.repo.SelectMultiByUids(uids)
}

// GetByUsernameAndPassword returns a user based on its username and passowrd,
// used for authentication.
func (s *userService) GetByNicknameAndPassword(nickname string, password string) (user dm.User, found bool) {
	if nickname == "" || password == "" {
		return dm.User{}, false
	}

	user, found = s.repo.SelectByNickname(nickname)

	//验证密码是否正确
	if ok, _ := dm.ValidatePassword(password, user.Password); ok {
		return user, true
	}
	return dm.User{}, false
}

//// Update updates every field from an existing User,
//// it's not safe to be used via public API,
//// however we will use it on the web/controllers/user_controller.go#PutBy
//// in order to show you how it works.
//func (s *userService) Update(id int64, user datamodels.User) (datamodels.User, error) {
//	user.ID = id
//	return s.repo.InsertOrUpdate(user)
//}

//// UpdatePassword updates a user's password.
//func (s *userService) UpdatePassword(id int64, newPassword string) (datamodels.User, error) {
//	// update the user and return it.
//	hashed, err := datamodels.GeneratePassword(newPassword)
//	if err != nil {
//		return datamodels.User{}, err
//	}

//	return s.Update(id, datamodels.User{
//		HashedPassword: hashed,
//	})
//}

//// UpdateUsername updates a user's username.
//func (s *userService) UpdateUsername(id int64, newUsername string) (datamodels.User, error) {
//	return s.Update(id, datamodels.User{
//		Username: newUsername,
//	})
//}

//// DeleteByID deletes a user by its id.
////
//// Returns true if deleted otherwise false.
//func (s *userService) DeleteByID(id int64) bool {
//	return s.repo.Delete(func(m datamodels.User) bool {
//		return m.ID == id
//	}, 1)
//}
