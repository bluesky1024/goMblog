package datamodels

type KafkaMsg struct {
	MsgId     int64
	Topic     string
	Partition int64
	Data      interface{}
}
