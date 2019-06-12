package relationService

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func newKafkaProducer(addr []string) (producer sarama.AsyncProducer, err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_2_0_0
	producer, err = sarama.NewAsyncProducer(addr, config)
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	return producer, nil
}

func (s *relationService) sendFollowMsg(msg dm.FollowMsg) {
	//kafkaConfig := conf.InitConfig("kafkaConfig.relation")

	msgStr, _ := json.Marshal(&msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic: "relationFollow",
		Key:   sarama.StringEncoder(msg.Uid),
		Value: sarama.ByteEncoder(msgStr),
	}

	s.kafkaProducer.Input() <- kafkaMsg
	select {
	case suc := <-s.kafkaProducer.Successes():
		fmt.Println("offset:", suc.Offset, "timestamp:", suc.Timestamp)
	case fail := <-s.kafkaProducer.Errors():
		fmt.Println("err:", fail.Err.Error())
	}
}

func (s *relationService) sendUnFollowMsg(msg dm.FollowMsg) {
	//kafkaConfig := conf.InitConfig("kafkaConfig.relation")

	msgStr, _ := json.Marshal(&msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic: "relationUnFollow",
		Key:   sarama.StringEncoder(msg.Uid),
		Value: sarama.ByteEncoder(msgStr),
	}

	s.kafkaProducer.Input() <- kafkaMsg
	select {
	case suc := <-s.kafkaProducer.Successes():
		fmt.Println("offset:", suc.Offset, "timestamp:", suc.Timestamp)
	case fail := <-s.kafkaProducer.Errors():
		fmt.Println("err:", fail.Err.Error())
	}
}

func (s *relationService) sendGroupMsg(msg dm.SetGroupMsg) {
	msgStr, _ := json.Marshal(&msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic: "relationSetGroup",
		Key:   sarama.StringEncoder(msg.Uid),
		Value: sarama.ByteEncoder(msgStr),
	}

	s.kafkaProducer.Input() <- kafkaMsg
	select {
	case suc := <-s.kafkaProducer.Successes():
		fmt.Println("offset:", suc.Offset, "timestamp:", suc.Timestamp)
	case fail := <-s.kafkaProducer.Errors():
		fmt.Println("err:", fail.Err.Error())
	}
}

func (s *relationService) HandleFollowMsg(msg dm.FollowMsg) (err error) {
	//新增粉丝表记录
	s.repo.AddOrUpdateFan(msg.FollowUid, msg.Uid)

	//维护粉丝数

	//维护关注数

	return err
}

func (s *relationService) HandleUnFollowMsg(msg dm.FollowMsg) (err error) {
	//删除粉丝表记录
	s.repo.DeleteFan(msg.FollowUid, msg.Uid)

	//维护粉丝数

	//维护关注数

	return err
}
