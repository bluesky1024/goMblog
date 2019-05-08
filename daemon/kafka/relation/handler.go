package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/daemon/kafka"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func newRelationHandler() (handler *kafkaConsumer.ConsumerGroupHandlerC) {
	handler = &kafkaConsumer.ConsumerGroupHandlerC{}

	//处理关注消息
	handler.RegisterHandler("relationFollow", handleFollow)

	////处理取关消息
	handler.RegisterHandler("relationUnFollow", handleUnFollow)

	return handler
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
