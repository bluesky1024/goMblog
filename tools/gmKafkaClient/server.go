package gmKafkaClient

type GmKafkaHandler interface {

}

var logType = "gmKafkaClient"

type Msg struct {
	Topic string
	Msg string
}