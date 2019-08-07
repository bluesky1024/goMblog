package relationService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	ds "github.com/bluesky1024/goMblog/datasource/dbSource"
	"github.com/bluesky1024/goMblog/datasource/redisSource"
	"github.com/bluesky1024/goMblog/repositories/dbRepo/relation"
	"github.com/bluesky1024/goMblog/repositories/redisRepo/relation"
	"github.com/bluesky1024/goMblog/tools/logger"

	"github.com/Shopify/sarama"
)

var logType = "relationService"

type RelationServicer interface {
	GetFollowsByUid(uid int64, page int, pageSize int) (follows []dm.FollowInfo, cnt int64)
	GetFansByUid(uid int64, page int, pageSize int) (fans []dm.FanInfo, cnt int64)
	Follow(uid int64, uidFollow int64) bool
	UnFollow(uid int64, uidFollow int64) bool
	CheckFollow(uidA int64, uidB int64) int
	CheckRelation(uidA int64, uidB int64) int8

	GetFollowCntByUids(uids []int64) (followCntMap map[int64]int64)
	GetFanCntByUids(uids []int64) (fanbCntMap map[int64]int64)

	/*分组管理*/
	GetGroupsByUid(uid int64) (groups []dm.FollowGroup, cnt int64)
	AddGroup(uid int64, groupName string) bool
	DelGroup(uid int64, groupId int64) bool
	UpdateGroup(group dm.FollowGroup) bool
	SetFollowGroup(uid int64, uidFollow int64, groupId int64) bool

	/*kafka关注取关分组管理补充操作*/
	HandleFollowMsg(msg dm.FollowMsg) (err error)
	HandleUnFollowMsg(msg dm.FollowMsg) (err error)

	ReleaseSrv() error
}

type relationService struct {
	//msgParam struct {
	//	kafkaProducer sarama.AsyncProducer
	//	msgTopic      struct {
	//		relationTrans string
	//		groupTrans    string
	//	}
	//}
	kafkaProducer sarama.AsyncProducer
	repo          *relationDbRepo.RelationDbRepository
	rdRepo        *relationRdRepo.RelationRbRepository
}

func NewRelationServicer() (RelationServicer, error) {
	//relation服务仓库初始化
	//mysql
	relationSourceM, err := ds.LoadRelation(true)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	relationSourceS, err := ds.LoadRelation(false)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	relationRepo := relationDbRepo.NewRelationRepository(relationSourceM, relationSourceS)

	//redis cluster
	redisSource, err := redisSource.LoadRelationRdSour()
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	redisRepo := relationRdRepo.NewRelationRdRepo(redisSource)

	//kafka生产者初始化
	//kafkaConfig := conf.InitConfig("kafkaConfig.relation")
	kafkaConfig := make(map[string]string)
	kafkaConfig["host"] = "0.0.0.0"
	kafkaConfig["port"] = "9092"
	kafkaProducer, err := newKafkaProducer([]string{kafkaConfig["host"] + ":" + kafkaConfig["port"]})
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}

	return &relationService{
		repo:          relationRepo,
		rdRepo:        redisRepo,
		kafkaProducer: kafkaProducer,
	}, nil
}

func (s *relationService) ReleaseSrv() (err error) {
	if s.kafkaProducer != nil {
		err = s.kafkaProducer.Close()
	}
	return err
}
