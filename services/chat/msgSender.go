package chatService

import (
	"encoding/json"
	"errors"
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

func (s *chatService) sendRoomStartMsg(uid int64) error {
	//kafkaConfig := conf.InitConfig("kafkaConfig.relation")

	msg := dm.RoomStatusSwitchMsg{
		Uid:    uid,
		Status: dm.WorkStatusOn,
	}

	msgStr, _ := json.Marshal(&msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic: "roomStart",
		Key:   sarama.StringEncoder(msg.Uid),
		Value: sarama.ByteEncoder(msgStr),
	}

	s.kafkaProducer.Input() <- kafkaMsg
	select {
	case suc := <-s.kafkaProducer.Successes():
		fmt.Println("offset:", suc.Offset, "timestamp:", suc.Timestamp)
		return nil
	case fail := <-s.kafkaProducer.Errors():
		logger.Err(logType, "err:"+fail.Err.Error())
		return errors.New("send new barrage failed" + string(msgStr) + fail.Err.Error())
	}
}

func (s *chatService) sendRoomStopMsg(uid int64) error {
	//kafkaConfig := conf.InitConfig("kafkaConfig.relation")

	msg := dm.RoomStatusSwitchMsg{
		Uid:    uid,
		Status: dm.WorkStatusInvalid,
	}

	msgStr, _ := json.Marshal(&msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic: "roomStop",
		Key:   sarama.StringEncoder(msg.Uid),
		Value: sarama.ByteEncoder(msgStr),
	}

	s.kafkaProducer.Input() <- kafkaMsg
	select {
	case suc := <-s.kafkaProducer.Successes():
		fmt.Println("offset:", suc.Offset, "timestamp:", suc.Timestamp)
		return nil
	case fail := <-s.kafkaProducer.Errors():
		logger.Err(logType, "err:"+fail.Err.Error())
		return errors.New("send new barrage failed" + string(msgStr) + fail.Err.Error())
	}
}

func (s *chatService) sendNewBarrageToRoom(msg dm.ChatBarrageInfo) error {
	//kafkaConfig := conf.InitConfig("kafkaConfig.relation")
	msgStr, _ := json.Marshal(&msg)
	kafkaMsg := &sarama.ProducerMessage{
		Topic: "newBarrageToRoom",
		Key:   sarama.StringEncoder(msg.Uid),
		Value: sarama.ByteEncoder(msgStr),
	}

	s.kafkaProducer.Input() <- kafkaMsg
	select {
	case suc := <-s.kafkaProducer.Successes():
		fmt.Println("offset:", suc.Offset, "timestamp:", suc.Timestamp)
		return nil
	case fail := <-s.kafkaProducer.Errors():
		logger.Err(logType, "err:"+fail.Err.Error())
		return errors.New("send new barrage failed" + string(msgStr) + fail.Err.Error())
	}
}
