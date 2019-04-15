package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bluesky1024/goMblog/config"
	"github.com/bluesky1024/goMblog/tools/logger"
)

var logType = "kafkaConsumerRelation"

func main() {
	initServ()

	relationConfig := conf.InitConfig("kafkaConfig.relation")
	fmt.Println(relationConfig)
	//relationConfig := make(map[string]string)
	//relationConfig["host"] = "0.0.0.0"
	//relationConfig["port"] = "9092"
	//relationConfig["groupId"] = "relation"
	//relationConfig["topic"] = "relation_trans"

	config := sarama.NewConfig()
	config.Version = sarama.V2_2_0_0
	config.Consumer.Return.Errors = true

	// Start with a client
	client, err := sarama.NewClient([]string{relationConfig["host"] + ":" + relationConfig["port"]}, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(relationConfig["groupId"], client)
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

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := []string{relationConfig["topic"]}
		handler := relationConsumerGroupHandler{}

		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
	defer func() {
		resourceRecycle()
	}()
}

func resourceRecycle() {
	//服务释放
	if relationSrv != nil {
		relationSrv.ReleaseSrv()
	}
}
