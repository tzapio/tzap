package codegeneration

import (
	"os"
	"path"
	"strings"

	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
)

func MakeCodeExtReplacer(language, extensionIn, extensionOut, mission, task string) func(t *tzap.Tzap) *tzap.Tzap {
	return func(t *tzap.Tzap) *tzap.Tzap {
		filein := t.Data["filepath"].(string)
		if path.Ext(filein) != extensionIn {
			return t
		}
		fileout := strings.TrimSuffix(t.Data["filepath"].(string), extensionIn) + extensionOut

		return t.
			HijackTzap(&tzap.Tzap{Name: "MakeCodeGO"}).
			AddSystemMessage("### The overall mission: \n"+mission).
			AddUserMessage(
				"###",
				"TASK: "+task,
				"PLAN: Do not write any text because this file will be saved directly to "+fileout,
				"TASKFILE: "+filein,
				"OUTFILE: "+fileout,
				"OUTPUT: "+language,
				"### EXAMPLE:",
				"{"+language+" code}").
			AddSystemMessage(
				"####",
				"####file: "+filein+"\n",
				t.Data["content"].(string),
			).
			MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
				if _, err := os.Stat(fileout); err == nil {
					return t.AddSystemMessage(
						"####",
						"####file: "+fileout+"\n",
						t.Data["content"].(string),
					)
				}
				return t
			})

	}
}

type BasicRefactoringConfig struct {
	FileIn           string   `json:"filein,omitempty"`
	FileOut          string   `json:"fileout,omitempty"`
	Mission          string   `json:"mission,omitempty"`
	Task             string   `json:"task,omitempty"`
	Plan             string   `json:"plan,omitempty"`
	OutputFormat     string   `json:"outputformat,omitempty"`
	Example          string   `json:"example,omitempty"`
	InspirationFiles []string `json:"inspirationfiles,omitempty"`
}

func MakeCode(config BasicRefactoringConfig) func(t *tzap.Tzap) *tzap.Tzap {

	systemMessage := ""
	if config.Mission != "" {
		systemMessage = "#### The overall mission: " + config.Mission
	}

	inspirationString := ""
	if len(config.InspirationFiles) > 0 {
		inspirationString = "\n#### The following files have been chosen as inspiration files\n"
		for _, v := range config.InspirationFiles {
			if v != "" {
				inspirationString += "\n####file: " + v + "\n"
				inspirationString += util.ReadFileP(v) + "\n"
			}
		}
	}

	systemMessage += inspirationString

	outputStr := ""
	if config.OutputFormat != "" {
		outputStr = "OUTPUT: " + config.OutputFormat
	}
	exampleStr := ""
	if config.Example != "" {
		exampleStr = "#### EXAMPLE:\n" + config.Example
	}
	planStr := ""
	if config.Plan != "" {
		planStr = config.Plan
	} else {
		planStr = "Do not write any text because this file will be saved directly to " + config.FileOut
	}

	return func(t *tzap.Tzap) *tzap.Tzap {
		return t.
			HijackTzap(&tzap.Tzap{Name: "MakeCodeGO"}).
			AddSystemMessage(
				systemMessage,
				"####",
				"TASK: "+config.Task,
				"PLAN: "+planStr,
				"TASKFILE: "+config.FileIn,
				"OUTFILE: "+config.FileOut,
				outputStr,
				exampleStr).
			MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
				t = t.AddUserMessage(
					"####\nfile in: "+config.FileIn+"\n####\n",
					util.ReadFileP(config.FileIn),
				)
				if config.FileIn != config.FileOut && config.FileOut != "" {
					if _, err := os.Stat(config.FileOut); err == nil {
						return t.AddUserMessage(
							"####\nfile out: "+config.FileOut+"\n####\n",
							util.ReadFileP(config.FileOut),
						)
					}
				}
				return t
			}).
			RequestChatCompletion()

	}
}
