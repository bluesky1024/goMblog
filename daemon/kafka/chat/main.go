package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/tools/logger"
)

var (
	logType = "kafkaConsumerChat"
	groupId = "chatSrv"
)

func main() {
	initServ()
	defer func() {
		resourceRecycle()
	}()

	//relationConfig := conf.InitConfig("kafkaConfig.relation")
	//fmt.Println(relationConfig)
	chatConfig := make(map[string]string)
	chatConfig["host"] = "0.0.0.0"
	chatConfig["port"] = "9092"

	config := sarama.NewConfig()
	config.Version = sarama.V2_2_0_0
	config.Consumer.Return.Errors = true

	// Start with a client
	client, err := sarama.NewClient([]string{chatConfig["host"] + ":" + chatConfig["port"]}, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(groupId, client)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			logger.Err(logType, err.Error())
		}
	}()

	msgHandler := newChatHandler()
	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		err := group.Consume(ctx, msgHandler.Topics, msgHandler)
		if err != nil {
			panic(err)
		}
	}
}
