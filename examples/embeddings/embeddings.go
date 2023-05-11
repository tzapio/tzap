package main

import (
	"github.com/tzapio/tzap/cli/cmd/util"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/embed"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	tutil "github.com/tzapio/tzap/pkg/util"
)

func main() {
	tzap.
		NewWithConnector(
			tzapconnect.WithConfig(
				config.Configuration{
					MD5Rewrites: true,
				})).
		WorkTzap(func(t *tzap.Tzap) {
			files, err := tutil.ListFilesInDir("./")
			if err != nil {
				panic(err)
			}
			files = util.GetNonExcludedFiles(files)
			embed.OutputEmbeddingsToFile(t, files)
		})
}
