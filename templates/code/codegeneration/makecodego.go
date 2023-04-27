package codegeneration

import (
	"strings"

	"github.com/tzapio/tzap/pkg/tzap"
)

func MakeCodeGO(mission, task string) func(t *tzap.Tzap) *tzap.Tzap {
	return func(t *tzap.Tzap) *tzap.Tzap {
		filein := t.Data["filepath"].(string)
		fileout := strings.TrimSuffix(t.Data["filepath"].(string), ".go") + ".go"
		if strings.HasSuffix(filein, "_test.go") || strings.Contains(filein, "file2") || strings.Contains(filein, "interfaces.go") {
			return t
		}
		return t.
			HijackTzap(&tzap.Tzap{Name: "MakeCodeGO"}).
			AddSystemMessage("### The overall mission: \n"+mission).
			AddUserMessage(
				"###",
				"TASK: "+task,
				"PLAN: Do not write any text because this file will be saved directly to "+fileout,
				"TASKFILE: "+filein,
				"OUTFILE: "+fileout,
				"OUTPUT: golang",
				"### EXAMPLE:",
				"EXAMPLE:",
				"{golang code}").
			AddSystemMessage(
				"###",
				"//file: "+filein+"\n",
				t.Data["content"].(string),
			).
			LoadTaskOrRequestNewTaskMD5(fileout)

		/*"These are the current existing files:\n//existing file (content redacted)",
		strings.Join(files, "\n//existing file (content redacted)"),*/

	}

}
