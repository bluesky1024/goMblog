package mblogService

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

func (s *mblogService) sendNewMblogMsg(msg dm.MblogNewMsg) {
	//kafkaConfig := conf.InitConfig("kafkaConfig.relation")

	msgStr, _ := json.Marshal(&msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic: "mblogNew",
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

func (s *mblogService) sendUpdateMblogMsg(msg dm.FollowMsg) {
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
