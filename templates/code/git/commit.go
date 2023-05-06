package git

import (
	"fmt"
	"os/exec"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func GitDiff() types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap] {
	templateF := func(t *tzap.Tzap) *tzap.ErrorTzap {
		diff := exec.Command("git", "diff",
			"--staged",
			"--patch-with-raw",
			"--unified=2",
			"--color=never",
			"--no-renames",
			"--ignore-space-change",
			"--ignore-all-space",
			"--ignore-blank-lines",
		)
		out, err := diff.CombinedOutput()
		if err != nil {
			return t.ErrorTzap(fmt.Errorf("could not get diff: %v", err))
		}

		t.Data["git-diff"] = string(out)
		return t.ErrorTzap(nil)
	}
	return types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap]{
		Name:     "GitDiff",
		Template: templateF,
	}
}

func ValidateDiff() types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap]{
		Name: "ValidateDiff",
		Template: func(t *tzap.Tzap) *tzap.ErrorTzap {
			diff := t.Data["git-diff"].(string)

			if diff == "" {
				return t.ErrorTzap(fmt.Errorf("diff is empty. Stage files to continue"))
			}

			return t.ErrorTzap(nil)
		}}
}

func GitCommit() types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap]{
		Name: "GitCommit",
		Template: func(t *tzap.Tzap) *tzap.ErrorTzap {
			content := t.Data["content"].(string)

			cmd := exec.Command("git", "commit", "-m", content)
			if err := cmd.Run(); err != nil {
				return t.ErrorTzap(fmt.Errorf("could not git commit. Content: %s, Error: %s", content, err))
			}

			return t.ErrorTzap(nil)
		}}
}
