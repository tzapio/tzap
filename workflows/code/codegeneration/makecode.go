package codegeneration

import (
	"path"
	"strings"

	"github.com/tzapio/tzap/pkg/tzap"
)

func MakeCode(language, extensionIn, extensionOut, mission, task string) func(t *tzap.Tzap) *tzap.Tzap {
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
				"EXAMPLE:",
				"{"+language+" code}").
			AddSystemMessage(
				"###",
				"###file: "+filein+"\n",
				t.Data["content"].(string),
			).
			LoadCompletionOrRequestCompletionMD5(fileout)

		/*"These are the current existing files:\n//existing file (content redacted)",
		strings.Join(files, "\n//existing file (content redacted)"),*/

	}

}
