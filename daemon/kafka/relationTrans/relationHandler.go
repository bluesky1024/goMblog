package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

type relationConsumerGroupHandler struct{}

func (h relationConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (h relationConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h relationConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		var err error
		var msgEncode = new(dm.FollowKafkaStruct)
		err = json.Unmarshal(msg.Value, msgEncode)
		if err != nil{
			fmt.Println("json unmarshal err:",err.Error())
			continue
		}

		switch msgEncode.Status {
		case dm.FollowStatusNormal:
			err = handleFollow(*msgEncode)
			break
		case dm.FollowStatusDelete:
			err = handleUnFollow(*msgEncode)
			break
		default:
			break
		}
		if err != nil {
			logger.Err(logType, err.Error())
		}
	}
	return nil
}

func handleFollow(msg dm.FollowKafkaStruct) error {
	relationSrv.HandleFollowMsg(msg)
	return nil
}

func handleUnFollow(msg dm.FollowKafkaStruct) error {
	relationSrv.HandleUnFollowMsg(msg)
	return nil
}
