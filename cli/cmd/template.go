package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/internal/logging/tl"
)

func init() {

}

var templateCmd = &cobra.Command{
	Aliases: []string{"s"},
	Use:     "template <file>",
	Short:   "Alpha",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tl.Logger.Println("Cobra CLI Template start")
	},
}
