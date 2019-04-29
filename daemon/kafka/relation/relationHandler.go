package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/tools/logger"
)

type TopicHandler func(msg sarama.ConsumerMessage) error

type relationConsumerGroupHandler struct {
	//map[topic]handler
	HandlerMap map[string]TopicHandler
}

func (h relationConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (h relationConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h relationConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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

func (h *relationConsumerGroupHandler) registerHandler(topic string, handler TopicHandler) {
	if h.HandlerMap == nil {
		h.HandlerMap = make(map[string]TopicHandler, 1)
	}
	h.HandlerMap[topic] = handler
}
