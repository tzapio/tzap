package cmd

import (
	"path"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/tzap"
)

func init() {
	RootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Aliases: []string{"i"},
	Use:     "install <name> <zip url>",
	Short:   "ALPHA: Install git packages",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tl.Logger.Println("Cobra CLI Install start")
		var name project.ProjectName = project.ProjectName(args[0])
		var projectDir project.ProjectDir = project.ProjectDir(path.Join("./.tzap-data/", string(name)))
		var zipURL = args[1]

		err := tzap.HandlePanic(func() {
			projectDB := project.ProjectDB{}
			projectDB[name] = projectDir
			t, err := initializeTzap(projectDB)
			if err != nil {
				panic(err)
			}
			defer t.HandleShutdown()
			t.ApplyWorkflow(cliworkflows.IndexZipFilesAndEmbeddings(name, projectDir, zipURL, false, false))
		})

		if err != nil {
			panic(err)
		}
	},
}
