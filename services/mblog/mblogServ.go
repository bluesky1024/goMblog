package mblogService

import (
	"github.com/Shopify/sarama"
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
	GetNormalByUid(uid int64, readAble []int8, page int, pageSize int) (mblogs []dm.MblogInfo, cnt int64)

	//	GetAll() []datamodels.User
	//	GetByID(id int64) (datamodels.User, bool)
	//	GetByUsernameAndPassword(username, userPassword string) (datamodels.User, bool)
	//	DeleteByID(id int64) bool

	//	Update(id int64, user datamodels.User) (datamodels.User, error)
	//	UpdatePassword(id int64, newPassword string) (datamodels.User, error)
	//	UpdateUsername(id int64, newUsername string) (datamodels.User, error)
}

type mblogService struct {
	kafkaProducer sarama.AsyncProducer
	repo          *mblogDbRepo.MblogDbRepository
}

// NewUserService returns the default user service.
func NewMblogServicer() (s MblogServicer, err error) {
	//id生成池初始化
	idGen.InitMidPool(10)

	//user服务仓库初始化
	mblogSourceM, err := ds.LoadMblogSour(true)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	mblogSourceS, err := ds.LoadMblogSour(false)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	mblogRepo := mblogDbRepo.NewMblogRepository(mblogSourceM, mblogSourceS)

	//kafka生产者初始化
	kafkaConfig := make(map[string]string)
	kafkaConfig["host"] = "0.0.0.0"
	kafkaConfig["port"] = "9092"
	kafkaProducer, err := newKafkaProducer([]string{kafkaConfig["host"] + ":" + kafkaConfig["port"]})
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}

	return &mblogService{
		repo:          mblogRepo,
		kafkaProducer: kafkaProducer,
	}, nil
}
