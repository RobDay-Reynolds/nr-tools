package main

import (
	"fmt"

	"github.com/monkeyherder/nr-tools/kafkaquota"
)

func main() {
	clients := kafkaquota.GetAllClients()

	quotas := kafkaquota.GenerateQuotasForClients(clients, 4)

	fmt.Println(fmt.Sprintf("%v", quotas))
}
