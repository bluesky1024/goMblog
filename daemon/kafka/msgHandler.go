package kafkaConsumer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/tools/logger"
)

var logType = "msg handle"

type TopicHandler func(msg sarama.ConsumerMessage) error

type ConsumerGroupHandler struct {
	Topics []string
	//map[topic]handler
	HandlerMap map[string]TopicHandler
}

func (h ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)

		tempHandler, ok := h.HandlerMap[msg.Topic]
		if !ok {
			logger.Err(logType, "no topic handler about "+msg.Topic)
			continue
		}

		err := tempHandler(*msg)
		if err != nil {
			logger.Err(logType, err.Error())
		}
	}
	return nil
}

func (h ConsumerGroupHandler) RegisterHandler(topic string, handler TopicHandler) {
	if h.Topics == nil {
		h.Topics = make([]string, 1)
	}
	if h.HandlerMap == nil {
		h.HandlerMap = make(map[string]TopicHandler, 1)
	}
	h.Topics = append(h.Topics, topic)
	h.HandlerMap[topic] = handler
}
