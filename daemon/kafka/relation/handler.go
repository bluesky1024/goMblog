package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func newRelationHandler() (topics []string, handler relationConsumerGroupHandler) {
	handler = relationConsumerGroupHandler{
		handlerMap: make(map[string]topicHandler, 1),
	}

	var topic string
	//处理关注消息
	topic = "relationFollow"
	topics = append(topics, topic)
	handler.registerHandler(topic, handleFollow)

	//处理取关消息
	topic = "relationUnFollow"
	topics = append(topics, topic)
	handler.registerHandler(topic, handleUnFollow)

	return topics, handler
}

func handleFollow(msg sarama.ConsumerMessage) (err error) {
	realMsg := new(dm.FollowMsg)
	err = json.Unmarshal(msg.Value, realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	err = relationSrv.HandleFollowMsg(*realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}

func handleUnFollow(msg sarama.ConsumerMessage) (err error) {
	realMsg := new(dm.FollowMsg)
	err = json.Unmarshal(msg.Value, realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	err = relationSrv.HandleUnFollowMsg(*realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}
