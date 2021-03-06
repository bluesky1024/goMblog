package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/daemon/kafka"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func newFeedHandler() (handler *kafkaConsumer.ConsumerGroupHandlerC) {
	handler = &kafkaConsumer.ConsumerGroupHandlerC{}

	//relationSrv
	//处理关注消息
	handler.RegisterHandler("relationFollow", handleFollow)
	//处理取关消息
	handler.RegisterHandler("relationUnFollow", handleUnFollow)
	//处理设置分组消息
	handler.RegisterHandler("relationSetGroup", handleSetGroup)

	//mblogSrv
	//处理用户发布新微博消息
	handler.RegisterHandler("mblogNew", handleMblogNew)

	return handler
}

func handleFollow(msg sarama.ConsumerMessage) (err error) {
	realMsg := new(dm.FollowMsg)
	err = json.Unmarshal(msg.Value, realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	err = feedSrv.HandleFollowMsg(*realMsg)
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
	err = feedSrv.HandleUnFollowMsg(*realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}

func handleSetGroup(msg sarama.ConsumerMessage) (err error) {
	realMsg := new(dm.SetGroupMsg)
	err = json.Unmarshal(msg.Value, realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	err = feedSrv.HandleSetGroupMsg(*realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}

func handleMblogNew(msg sarama.ConsumerMessage) (err error) {
	realMsg := new(dm.MblogNewMsg)
	err = json.Unmarshal(msg.Value, realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	err = feedSrv.HandleMblogNewMsg(*realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}
