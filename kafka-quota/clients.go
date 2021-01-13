package kafkaquota

type KafkaClient struct {
	ID                string `json:"clientId"`
	AverageThroughput int
	MaxThroughput     int
}
