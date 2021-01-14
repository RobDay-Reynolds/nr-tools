package main

import (
	"html/template"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/monkeyherder/nr-tools/kafkaquota"
)

func main() {
	var nrAPIKey string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "nr-api-key",
				Value:       "",
				Usage:       "New Relic User API Key",
				Destination: &nrAPIKey,
			},
		},

		Name:  "kafka-quotas",
		Usage: "Calculate kafka quotas for existing producers",
		Action: func(c *cli.Context) error {
			clients := kafkaquota.GetAllClients(nrAPIKey)

			quotas := kafkaquota.GenerateQuotasForClients(clients, 4)
			sort.Slice(quotas, func(i, j int) bool {
				return quotas[i].ClientID < quotas[j].ClientID
			})

			tmpl, _ := template.New("quotas").Parse(kafkaquota.QuotasTemplate)
			tmpl.Execute(os.Stdout, quotas)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
