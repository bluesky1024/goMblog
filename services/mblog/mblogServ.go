package mblogService

import (
	"errors"
	"fmt"
	dm "github.com/bluesky1024/goMblog/datamodels"
	ds "github.com/bluesky1024/goMblog/datasource/dbSource"
	"github.com/bluesky1024/goMblog/repositories/dbRepo/mblog"
	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/tools/logger"
)

var logType = "userService"

type MblogServicer interface {
	Create(uid int64, content string, readAble int8, originUid int64, originMid int64) (mblog dm.MblogInfo, err error)
	GetByMid(mid int64) (mblog dm.MblogInfo, found bool)
	GetMultiByMids(mids []int64) map[int64]dm.MblogInfo
	GetByUid(uid int64, page int, pageSize int) (mblogs map[int64]dm.MblogInfo)

	//	GetAll() []datamodels.User
	//	GetByID(id int64) (datamodels.User, bool)
	//	GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool)
	//	DeleteByID(id int64) bool

	//	Update(id int64, user datamodels.User) (datamodels.User, error)
	//	UpdatePassword(id int64, newPassword string) (datamodels.User, error)
	//	UpdateUsername(id int64, newUsername string) (datamodels.User, error)
}

type mblogService struct {
	repo *mblogDbRepo.MblogDbRepository
}

// NewUserService returns the default user service.
func NewMblogServicer() MblogServicer {
	//id生成池初始化
	idGen.InitMidPool(10)

	//user服务仓库初始化
	mblogSourceM, err := ds.LoadMblogSour(true)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	mblogSourceS, err := ds.LoadMblogSour(false)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	mblogRepo := mblogDbRepo.NewMblogRepository(mblogSourceM, mblogSourceS)

	return &mblogService{
		repo: mblogRepo,
	}
}

func (m *mblogService) Create(uid int64, content string, readAble int8, originUid int64, originMid int64) (mblog dm.MblogInfo, err error) {
	fmt.Println("enter create")
	if uid <= 0 || content == "" || len(content) > 140 {
		return dm.MblogInfo{}, errors.New("invalid mblog info")
	}

	insertData := dm.MblogInfo{
		Uid:       uid,
		Content:   content,
		OriginMid: originMid,
		OriginUid: originUid,
		ReadAble:  readAble,
		Status:    dm.MblogStatusNormal,
	}
	//生成专属mid
	insertData.Mid, err = idGen.GenMidId()
	if err != nil {
		return dm.MblogInfo{}, err
	}

	affected, err := m.repo.Insert(insertData)
	if err != nil || affected == 0 {
		logger.Err(logType, err.Error())
		return dm.MblogInfo{}, err
	}
	return insertData, nil
}

func (m *mblogService) GetByMid(mid int64) (mblog dm.MblogInfo, found bool) {
	if mid <= 0 {
		return dm.MblogInfo{}, false
	}
	return m.repo.SelectByMid(mid)
}

func (m *mblogService) GetMultiByMids(mids []int64) map[int64]dm.MblogInfo {
	res := m.repo.SelectMultiByMids(mids)
	return res
}

func (m *mblogService) GetByUid(uid int64, page int, pageSize int) (mblogs map[int64]dm.MblogInfo) {
	if uid <= 0 {
		return nil
	}
	mblogs = m.repo.SelectByUid(uid, page, pageSize)
	return mblogs
}
