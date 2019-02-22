package gmKafkaClient

import (
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/tools/logger"
)
func NewGMKafkaProducer(addr []string) sarama.AsyncProducer{
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_2_0_0
	producer, err := sarama.NewAsyncProducer(addr, config)
	if err != nil {
		logger.Err(logType,err.Error())
	}
	return producer
}