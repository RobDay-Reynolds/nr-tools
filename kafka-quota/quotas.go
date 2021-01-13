package kafkaquota

type KafkaQuota struct {
	ClientID string `json:"clientId"`
	Quota    int    `json:"perBrokerProducerByteRate"`
}

func GenerateQuotasForClients(clients []KafkaClient, throughputMultiplier int) []KafkaQuota {
	quotas := []KafkaQuota{}

	for _, client := range clients {
		quota := KafkaQuota{
			ClientID: client.ID,
		}

		roundingMultiple := 0

		switch throughput := client.MaxThroughput; {
		case throughput > 1048576:
			roundingMultiple = 1048576
		case throughput > 1024:
			roundingMultiple = 1024
		}

		quota.Quota = roundToMultiple(client.MaxThroughput*5, roundingMultiple)

		quotas = append(quotas, quota)
	}

	return quotas
}

func roundToMultiple(numToRound, multiple int) int {
	if multiple == 0 {
		return numToRound
	}

	remainder := numToRound % multiple
	if remainder == 0 {
		return numToRound
	}

	return numToRound + multiple - remainder
}