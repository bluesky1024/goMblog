package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/daemon/kafka"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func newChatHandler() (handler *kafkaConsumer.ConsumerGroupHandlerC) {
	handler = &kafkaConsumer.ConsumerGroupHandlerC{}

	//处理房间开关消息
	handler.RegisterHandler("roomStart", handleRoomStart)
	handler.RegisterHandler("roomStop", handleRoomStop)

	//处理新消息推送至房间
	handler.RegisterHandler("newBarrageToRoom", handleNewMsgToRoom)

	return handler
}

func handleRoomStart(msg sarama.ConsumerMessage) (err error) {
	realMsg := new(dm.RoomStatusSwitchMsg)
	err = json.Unmarshal(msg.Value, realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	err = chatSrv.HandleRoomStartMsg(*realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}

func handleRoomStop(msg sarama.ConsumerMessage) (err error) {
	realMsg := new(dm.RoomStatusSwitchMsg)
	err = json.Unmarshal(msg.Value, realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	err = chatSrv.HandleRoomStopMsg(*realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}

func handleNewMsgToRoom(msg sarama.ConsumerMessage) (err error) {
	realMsg := new(dm.ChatBarrageInfo)
	err = json.Unmarshal(msg.Value, realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	err = chatSrv.HandleNewBarrageToRoomMsg(*realMsg)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}
