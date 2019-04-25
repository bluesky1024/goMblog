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
	handler.registerHandler("relation_follow", handleFollow)
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
