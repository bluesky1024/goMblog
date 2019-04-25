package kafkaClient

import (
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/tools/logger"
	"sync"
)

type KafkaProducer interface {
	//异步推送消息，无返回结果
	SendMsg(msg KafkaMsg)
	//推送消息，确认推送结果后才返回
	SendMsgGetRes(msg KafkaMsg) (err error)
}

type kafkaProducer struct {
	mutex sync.RWMutex

	producer sarama.AsyncProducer
}

func NewKafkaProducer(addr []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_2_0_0
	producer, err := sarama.NewAsyncProducer(addr, config)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return producer
}
