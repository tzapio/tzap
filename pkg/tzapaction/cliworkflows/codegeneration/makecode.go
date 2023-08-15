package codegeneration

import (
	"os"

	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/util"
)

func MakeCode(args *actionpb.RefactorArgs) func(t *tzap.Tzap) *tzap.Tzap {
	systemMessage := ""
	if args.Mission != "" {
		systemMessage = "#### The overall mission: " + args.Mission
	}

	inspirationString := ""
	added := false
	if len(args.InspirationFiles) > 0 {
		for _, v := range args.InspirationFiles {
			if v != "" {
				if _, err := os.Stat(v); err == nil {
					inspirationString += "\n####file: " + v + "\n"
					inspirationString += util.ReadFileP(v) + "\n"

					added = true
				}
			}
		}
	}
	if added {
		inspirationString = "\n####The following files have been chosen as inspiration files\n" + inspirationString
	}

	systemMessage += inspirationString

	outputStr := ""
	if args.OutputFormat != "" {
		outputStr = "OUTPUT: " + args.OutputFormat
	}
	exampleStr := ""
	if args.Example != "" {
		exampleStr = "####EXAMPLE:\n" + args.Example
	}
	planStr := ""
	if args.Plan != "" {
		planStr = args.Plan
	} else {
		planStr = "Do not write any text because this file will be saved directly to " + args.FileOut
	}
	if args.FileOut == "" {
		args.FileOut = args.FileIn
	}
	return func(t *tzap.Tzap) *tzap.Tzap {
		return t.
			HijackTzap(&tzap.Tzap{Name: "MakeCodeGO"}).
			AddSystemMessage(
				systemMessage,
				"Do not write anything using #### as that is the delimeter. The first line of your output must be the file content and nothing else.",
				"####",
				"TASK: "+args.Task,
				"PLAN: "+planStr,
				"TASKFILE: "+args.FileIn,
				"OUTFILE: "+args.FileOut,
				outputStr,
				exampleStr).
			MutationTzap(func(t *tzap.Tzap) *tzap.Tzap {
				if args.FileIn != "" {
					if _, err := os.Stat(args.FileIn); err == nil {
						t = t.AddSystemMessage(
							"####file in: "+args.FileIn,
							util.ReadFileP(args.FileIn),
						)
					}
				}

				if args.FileIn != args.FileOut && args.FileOut != "" {
					if _, err := os.Stat(args.FileOut); err == nil {
						return t.
							AddSystemMessage(
								"####file out: "+args.FileOut,
								util.ReadFileP(args.FileOut),
							)
					}
				}
				return t
			}).
			RequestChatCompletion()

	}
}
