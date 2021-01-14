package kafkaquota

import "fmt"

const QuotaLimit = 62914560
const QuotasTemplate = `{{ range . }}
  - clientId: {{ .ClientID }}
    perBrokerProducerByteRate: {{ .Quota }} # {{ .HumanReadableQuota }}
{{ end }}`

type KafkaQuota struct {
	ClientID           string `json:"clientId"`
	Quota              int    `json:"perBrokerProducerByteRate"`
	HumanReadableQuota string
}

func GenerateQuotasForClients(clients []KafkaClient, throughputMultiplier int) []KafkaQuota {
	quotas := []KafkaQuota{}

	for _, client := range clients {
		quota := KafkaQuota{
			ClientID: client.ID,
		}

		roundingMultiple := 1
		ratePrefix := ""

		finalMultiplier := throughputMultiplier

		clientThroughput := int(client.AverageThroughput)
		if float64(client.MaxThroughput) > client.AverageThroughput*2 {
			clientThroughput = client.MaxThroughput
			finalMultiplier = throughputMultiplier / 2
		}

		if clientThroughput > 1048576 {
			roundingMultiple = 1048576
			ratePrefix = "M"
		} else if clientThroughput > 1024 {
			roundingMultiple = 1024
			ratePrefix = "k"
		}

		limit := roundToMultiple(clientThroughput*finalMultiplier, roundingMultiple)
		quota.Quota = numberOrCeiling(limit)
		quota.HumanReadableQuota = fmt.Sprintf("%d %sB/sec", quota.Quota/roundingMultiple, ratePrefix)

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

func numberOrCeiling(number int) int {
	if number > QuotaLimit {
		return QuotaLimit
	}

	return number
}
