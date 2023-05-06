package translate

import (
	"strings"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
)

func TranslateCodeFromTo(from, to, outDir, mission, task string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "translateCodeFromTo",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			filein := t.Data["filepath"].(string)
			fileout := strings.TrimSuffix(t.Data["filepath"].(string), ".go") + ".ts"
			fileout = strings.ReplaceAll(fileout, "/tzapio/tzap/", "/tzapio/tzap/ts/src/")
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
					"{"+to+" code}").
				AddSystemMessage(
					"###",
					"//file: "+filein+"\n",
					t.Data["content"].(string),
				).
				LoadTaskOrRequestNewTaskMD5(fileout)

		}}

}

func MakeCodeTSMessage(mission, task, filein, fileout string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "makeCodeTSMessage",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				AddSystemMessage("### The overall mission: \n"+mission).
				AddUserMessage(
					"###",
					"TASK: "+task,
					"PLAN: Do not write any text because this file will be saved directly to "+fileout,
					"TASKFILE: "+filein,
					"OUTFILE: "+fileout,
					"OUTPUT: ts",
					"### EXAMPLE:",
					"EXAMPLE:",
					"{ts code}").
				AddSystemMessage(
					"###",
					"//file: "+filein+"\n",
					util.ReadFileP(filein),
				).
				LoadTaskOrRequestNewTask(fileout)
		},
	}
}
