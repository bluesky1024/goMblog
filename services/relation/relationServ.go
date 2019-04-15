package relationService

import (
	"encoding/json"
	"fmt"
	"github.com/bluesky1024/goMblog/config"
	dm "github.com/bluesky1024/goMblog/datamodels"
	ds "github.com/bluesky1024/goMblog/datasource/dbSource"
	"github.com/bluesky1024/goMblog/repositories/dbRepo/relation"
	"github.com/bluesky1024/goMblog/tools/gmKafkaClient"
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

	/*分组信息*/
	GetGroupsByUid(uid int64) (groups []dm.FollowGroup, cnt int64)
	AddGroup(uid int64, groupName string) bool
	DelGroup(uid int64, groupId int64) bool
	UpdateGroup(group dm.FollowGroup) bool
	SetFollowGroup(uid int64, uidFollow int64, groupId int32) bool

	/*kafka关注取关分组管理补充操作*/
	HandleFollowMsg(msg dm.FollowKafkaStruct) (err error)
	HandleUnFollowMsg(msg dm.FollowKafkaStruct) (err error)

	ReleaseSrv() error
}

type relationService struct {
	kafkaProducer sarama.AsyncProducer
	repo          *relationDbRepo.RelationDbRepository
}

func NewRelationServicer() RelationServicer {
	//relation服务仓库初始化
	relationSourceM, err := ds.LoadRelation(true)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	relationSourceS, err := ds.LoadRelation(false)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	relationRepo := relationDbRepo.NewRelationRepository(relationSourceM, relationSourceS)

	//kafka生产者初始化
	kafkaConfig := conf.InitConfig("kafkaConfig.relation")
	kafkaProducer := gmKafkaClient.NewGMKafkaProducer([]string{kafkaConfig["host"] + ":" + kafkaConfig["port"]})

	return &relationService{
		repo:          relationRepo,
		kafkaProducer: kafkaProducer,
	}
}

func (s *relationService) GetFollowsByUid(uid int64, page int, pageSize int) (follows []dm.FollowInfo, cnt int64) {
	return s.repo.SelectMultiFollowsByUid(uid, page, pageSize)
}

func (s *relationService) GetFansByUid(uid int64, page int, pageSize int) (fans []dm.FanInfo, cnt int64) {
	return s.repo.SelectMultiFansByUid(uid, page, pageSize)
}

func (s *relationService) GetGroupsByUid(uid int64) (groups []dm.FollowGroup, cnt int64) {
	return s.repo.SelectMultiGroupsByUid(uid)
}

func (s *relationService) AddGroup(uid int64, groupName string) bool {
	return s.repo.AddOrUpdateGroup(uid, groupName)
}

func (s *relationService) DelGroup(uid int64, groupId int64) bool {
	return s.repo.DeleteGroupByUidAndGroupId(uid, groupId)
}

func (s *relationService) UpdateGroup(group dm.FollowGroup) bool {
	return s.repo.UpdateGroupById(group)
}

func (s *relationService) Follow(uid int64, uidFollow int64) bool {
	//修改follow表
	succ := s.repo.AddOrUpdateFollow(uid, uidFollow)

	//其他操作加入消息队列进行操作
	if succ {
		msg := dm.FollowKafkaStruct{
			Uid:       uid,
			FollowUid: uidFollow,
			Status:    dm.FollowStatusNormal,
		}
		s.sendKafkaMsg(msg)
	}
	return succ
}

func (s *relationService) SetFollowGroup(uid int64, uidFollow int64, groupId int32) bool {
	return false
}

func (s *relationService) UnFollow(uid int64, uidFollow int64) bool {
	//修改follow表
	succ := s.repo.DeleteFollow(uid, uidFollow)

	//其他操作加入消息队列进行操作
	if succ {
		msg := dm.FollowKafkaStruct{
			Uid:       uid,
			FollowUid: uidFollow,
			Status:    dm.FollowStatusDelete,
		}
		s.sendKafkaMsg(msg)
	}
	return succ
}

func (s *relationService) CheckFollow(uidA int64, uidB int64) int {
	if uidA == 0 || uidB == 0 || uidA == uidB {
		return 0
	}

	info, found := s.repo.SelectFollowByUid(uidA, uidB)
	if !found || info.Status == dm.FollowStatusDelete {
		return 0
	}

	return 1
}

func (s *relationService) CheckRelation(uidA int64, uidB int64) int8 {
	if uidA == 0 || uidB == 0 {
		return dm.RelationNone
	}
	if uidA == uidB {
		return dm.RelationSelf
	}
	return dm.RelationNone
}

func (s *relationService) sendKafkaMsg(msg dm.FollowKafkaStruct) {
	kafkaConfig := conf.InitConfig("kafkaConfig.relation")
	msgStr, _ := json.Marshal(&msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic:     kafkaConfig["topic"],
		Key:       sarama.StringEncoder(kafkaConfig["topic"]),
		Partition: 1,
		Value:     sarama.ByteEncoder(msgStr),
	}

	s.kafkaProducer.Input() <- kafkaMsg
	select {
	case suc := <-s.kafkaProducer.Successes():
		fmt.Println("offset:", suc.Offset, "timestamp:", suc.Timestamp)
	case fail := <-s.kafkaProducer.Errors():
		fmt.Println("err:", fail.Err.Error())
	}
}

func (s *relationService) HandleFollowMsg(msg dm.FollowKafkaStruct) (err error) {
	//新增粉丝表记录
	s.repo.AddOrUpdateFan(msg.FollowUid, msg.Uid)

	return err
}

func (s *relationService) HandleUnFollowMsg(msg dm.FollowKafkaStruct) (err error) {
	//删除粉丝表记录
	s.repo.DeleteFan(msg.FollowUid, msg.Uid)

	return err
}

func (s *relationService) ReleaseSrv() (err error) {
	if s.kafkaProducer != nil {
		err = s.kafkaProducer.Close()
	}
	return err
}
