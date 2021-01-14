package kafkaquota

import (
	"fmt"
	"log"
	"os"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"
)

type KafkaClient struct {
	ID                string `json:"clientId"`
	AverageThroughput float64
	MaxThroughput     int
}

func GetAllClients() []KafkaClient {
	// Initialize the client configuration.  A Personal API key is required to
	// communicate with the backend API.
	cfg := config.New()
	cfg.PersonalAPIKey = os.Getenv("NEW_RELIC_API_KEY")

	// Initialize the client.
	client := nerdgraph.New(cfg)

	// Execute a NRQL query to retrieve the average duration of transactions for
	// the "Example application" app.
	query := `
	query($accountId: Int!, $nrqlQuery: Nrql!) {
		actor {
			account(id: $accountId) {
				nrql(query: $nrqlQuery, timeout: 5) {
					results
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"accountId": 1,
		"nrqlQuery": "SELECT average(produceByteRateSec), max(produceByteRateSec) FROM KafkaClientStats FACET clientId WHERE clientId LIKE 'producer.%' AND produceByteRateSec > 1048576 SINCE 2 WEEKS AGO",
	}

	resp, err := client.Query(query, variables)
	if err != nil {
		log.Fatal("error running NerdGraph query: ", err)
	}

	queryResp := resp.(nerdgraph.QueryResponse)
	actor := queryResp.Actor.(map[string]interface{})
	account := actor["account"].(map[string]interface{})
	nrql := account["nrql"].(map[string]interface{})
	results := nrql["results"].([]interface{})

	var clients []KafkaClient
	for _, r := range results {
		data := r.(map[string]interface{})
		clients = append(
			clients,
			KafkaClient{
				clientId:          data["clientId"].(string),
				averageThroughput: data["average.produceByteRateSec"].(float64),
				maxThroughput:     int(data["max.produceByteRateSec"].(float64)),
			},
		)
	}

	// Output the raw time series values for transaction duration.
	fmt.Printf(": %v\n", clients)
}
