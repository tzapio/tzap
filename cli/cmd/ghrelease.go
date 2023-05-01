package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/pkg/util/stdin"
)

var ghrelease = &cobra.Command{
	Use:   "ghrelease",
	Short: "Generate a gh release using ChatGPT",
	Long: `Prompts ChatGPT to generate a commit message and commits it to the current git repo.
	The generated commit message is based on the diff of the currently staged files.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get previous and current tags
		prevTag := ""
		currentTag := args[0]
		prevTagCmd := exec.Command("git", "describe", "--tags", "--abbrev=0", "HEAD^")
		prevTagOutput, err := prevTagCmd.CombinedOutput()
		if err == nil {
			prevTag = strings.TrimSpace(string(prevTagOutput))
		}

		// Get git commits from last tag
		commitsCmd := exec.Command("git", "log", "--pretty=format:%s", fmt.Sprintf("%s..HEAD", prevTag))
		commitsOutput, err := commitsCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Could not get git commits:", err)
			return
		}

		// Create title and summary of changes
		commits := strings.Split(strings.TrimSpace(string(commitsOutput)), "\n")
		title := fmt.Sprintf("Release %s", currentTag)
		summary := ""
		for _, commit := range commits {
			summary += fmt.Sprintf("* %s \n", commit)
		}
		println(summary)

		url, err := exec.Command("git", "ls-remote", "--get-url").Output()
		if err != nil {
			fmt.Println("Could not get remote URL:", err)
			return
		}

		t := tzap.
			NewWithConnector(tzapconnect.WithConfig(config.Configuration{SupressLogs: true, OpenAIModel: modelMap[settings.Model]})).
			SetHeader(fmt.Sprintf(`User will post a changelog. You will output the Github release. Include a compare tag url.

Template:
{"title":{title}, "notes":{github release in markdown}\}

Github repo: ` + string(url))).
			AddUserMessage(fmt.Sprintf("TITLE: %s\n\nGITCOMMITS:\n%s", title, summary))

		res := t.RequestChat()

		// Parse the JSON object
		var data map[string]string
		err = json.Unmarshal([]byte(res.Data["content"].(string)), &data)
		if err != nil {
			fmt.Println("Could not parse JSON object:", err)
			return
		}

		// Get title and notes from the JSON object
		notes := data["notes"]
		if !stdin.ConfirmToContinue() {
			return
		}
		// Create release
		releaseCmd := exec.Command("gh", "release", "create", currentTag, "--prerelease", "--title", title, "--notes", notes)
		releaseCmd.Stderr = os.Stderr
		if err := releaseCmd.Run(); err != nil {
			fmt.Printf("Could not create release. Title: %s, Notes: %s, Error: %s\n", title, notes, err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(ghrelease)

}
