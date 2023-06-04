package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resetting embeddings and other files",

	Run: func(cmd *cobra.Command, args []string) {
		// delete .tzap-data/embeddingsCache.db, fileembeddings.db, and filesTimestamps.db
		tzapDataFilesToDelete := []string{
			"embeddingsCache.db",
			"fileembeddings.db",
			"filesTimestamps.db",
		}

		for _, file := range tzapDataFilesToDelete {
			filePath := path.Join(".tzap-data/", file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				println("Ignored - File " + file + " does not exist")
				continue
			}

			err := os.Remove(filePath)
			if err != nil {
				println("Error deleting file " + filePath)
				continue
			}
			println("Deleted " + filePath)

		}

	},
}

func init() {
	RootCmd.AddCommand(ResetCmd)
}
