package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
)

func main() {
	openai_apikey, err := tzapconnect.LoadOPENAI_APIKEY()
	if err != nil {
		panic(err)
	}
	tzap.
		NewWithConnector(
			tzapconnect.WithConfig(openai_apikey, config.Configuration{MD5Rewrites: true})).
		WorkTzap(func(t *tzap.Tzap) {
			err := embed.CreateQueryJSON(t, "query.json", "Write a tzap to get embeddings from files in a directory")
			if err != nil {
				panic(err)
			}
		})
}
