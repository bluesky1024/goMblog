package userDbRepo

import (
	"errors"
	//	"fmt"

	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"github.com/go-xorm/xorm"
)

var logType string = "userDbRepo"

//type Query func(datamodels.User) bool

//type UserRepository interface {
//	//Exec(query Query, action Query, limit int, mode int) (ok bool)

//	Select(query Query) (user datamodels.User, found bool)
//	SelectMany(query Query, limit int) (results []datamodels.User)

//	InsertOrUpdate(user datamodels.User) (updatedUser datamodels.User, err error)
//	Delete(query Query, limit int) (deleted bool)
//}

type UserDbRepository struct {
	sourceM *xorm.Engine
	sourceS *xorm.Engine
}

func NewUserRepository(sourceM *xorm.Engine, sourceS *xorm.Engine) *UserDbRepository {
	return &UserDbRepository{
		sourceM: sourceM,
		sourceS: sourceS,
	}
}

func (r *UserDbRepository) Select(u *dm.User) (found bool) {
	found, err := r.sourceS.Get(u)
	if err != nil {
		logger.Err(logType, err.Error())
		return false
	}
	return found
}

func (r *UserDbRepository) SelectByNickname(nickname string) (user dm.User, found bool) {
	found, err := r.sourceS.Where("nick_name=?", nickname).Get(&user)
	if err != nil {
		logger.Err(logType, err.Error())
		return dm.User{}, false
	}
	return user, found
}

func (r *UserDbRepository) SelectByUid(uid int64) (user dm.User, found bool) {
	found, err := r.sourceS.Where("uid = ?", uid).Get(&user)
	if err != nil {
		logger.Err(logType, err.Error())
		return dm.User{}, false
	}
	return user, found
}

func (r *UserDbRepository) SelectMultiByUids(uids []int64) (users map[int64]dm.User, err error) {
	users = make(map[int64]dm.User)
	var userArr []dm.User
	err = r.sourceS.In("uid",uids).Find(&userArr)
	if err != nil{
		return users,err
	}

	for _,user := range userArr {
		users[user.Uid] = user
	}
	return users,nil
}

func (r *UserDbRepository) Insert(u dm.User) (affected int64, err error) {
	affected, err = r.sourceM.Insert(u)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return affected, err
}

//func (r *UserDbRepository) Update(u dm.User) (affected int, err error) {
//	affected, err := r.source.Update(u)
//	if err != nil {
//		logger.Err(logType, err)
//	}
//	return affected, err
//}

func (r *UserDbRepository) UpdateByUid(u dm.User) (affected int64, err error) {
	if u.Uid <= 0 {
		err = errors.New("update user whit invalid id")
		logger.Err(logType, err.Error())
		return 0, err
	}
	affected, err = r.sourceM.Where("uid=?", u.Uid).Cols("password").Update(&u)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return affected, err
}
