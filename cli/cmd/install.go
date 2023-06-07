package cmd

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/cli/cmd/cliworkflows"
	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator/cmdinstance"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/tzap"
)

func init() {
	RootCmd.AddCommand(installCmd)
}

func GetZipUrlFromGithubUrl(githubUrl string) (string, error) {
	parsed, err := url.Parse(githubUrl)
	if err != nil {
		return "", err
	}

	// Check if url domain is github.com
	if parsed.Host != "github.com" {
		return "", fmt.Errorf("unsupported domain found: %v", parsed.Host)
	}

	// Extract user and repo from url path
	path := strings.TrimSuffix(parsed.Path, ".git")
	parts := strings.Split(path, "/")[1:]
	if len(parts) < 2 {
		return "", fmt.Errorf("could not parse user and repo from repo URL")
	}

	// Build the zip URL
	repoPath := strings.Join(parts, "/")
	zipUrl := fmt.Sprintf("https://github.com/%s/archive/refs/heads/main.zip", repoPath)

	return zipUrl, nil
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

		zipURL, err := GetZipUrlFromGithubUrl(args[1])
		if err != nil {
			panic(err)
		}

		err = tzap.HandlePanic(func() {
			zipProject, err := cmdinstance.NewLocalZipProject(string(name), "/", zipURL, []string{}, []string{})
			if err != nil {
				panic(err)
			}

			t, err := initializeTzap()
			if err != nil {
				panic(err)
			}
			t = t.AddTzap(&tzap.Tzap{Name: "loadAndSearchEmbeddings"}).
				MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
					t.C = project.SetProjectInContext(t.C, zipProject)
					return t
				})
			defer t.HandleShutdown()
			t.ApplyWorkflow(cliworkflows.IndexZipFilesAndEmbeddings(name, projectDir, zipURL, false, false))
		})

		if err != nil {
			panic(err)
		}
	},
}
