package kafkaClient

type GmKafkaHandler interface {
}

var logType = "kafkaClient"

type KafkaMsg struct {
	MsgId     int64
	Topic     string
	Partition int64
	Data      interface{}
}
