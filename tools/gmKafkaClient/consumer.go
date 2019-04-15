package gmKafkaClient

import (
	"github.com/Shopify/sarama"
)

func NewGMKafkaConsumer(addr []string) sarama.Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_2_0_0
	consumer, _ := sarama.NewConsumer([]string{"10.222.76.152:9092"}, config)
	return consumer
}
