package gocode

import (
	"fmt"
	"strings"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

const (
	maxTokensForGPT4    = 8000
	maxTokensForDefault = 4000
)

func DeserializedArguments(dataname string, args []string) types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "deserializedArguments",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			t.Data[dataname] = strings.Join(args, " ")
			return t
		}}
}
func SetContextSize() types.NamedTemplate[*tzap.Tzap, *tzap.Tzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.Tzap]{
		Name: "setContextSize",
		Template: func(t *tzap.Tzap) *tzap.Tzap {
			settings := config.FromContext(t.C)
			var contextSize int
			if settings.OpenAIModel == "gpt4" {
				contextSize = maxTokensForGPT4
			} else {
				contextSize = maxTokensForDefault
			}
			t.Data["contextSize"] = contextSize

			return t
		}}
}

func CountTokens() types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap]{
		Name: "countTokens",
		Template: func(t *tzap.Tzap) *tzap.ErrorTzap {
			diff := t.Data["git-diff"].(string)
			headerCount, err := t.CountTokens(t.Parent.Header)
			if err != nil {
				return t.ErrorTzap(fmt.Errorf("could not count tokens: %v", err))
			}
			contentTokens, err := t.CountTokens(diff)
			if err != nil {
				return t.ErrorTzap(fmt.Errorf("could not count tokens: %v", err))
			}
			t.Data["headerTokens"] = headerCount
			t.Data["contentTokens"] = contentTokens

			return t.ErrorTzap(nil)
		}}
}

func TruncateTokens() types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap] {
	return types.NamedTemplate[*tzap.Tzap, *tzap.ErrorTzap]{
		Name: "truncateTokens",
		Template: func(t *tzap.Tzap) *tzap.ErrorTzap {
			contextSize := t.Data["contextSize"].(int)
			headerTokens := t.Data["headerTokens"].(int)
			contentTokens := t.Data["contentTokens"].(int)

			max := contextSize - headerTokens - 1500
			if contentTokens >= max {
				offsetStart := 0
				offsetEnd := 0 + max
				diff := t.Data["git-diff"].(string)
				truncatedDiff, err := t.OffsetTokens(diff, offsetStart, offsetEnd)
				if err != nil {
					return t.ErrorTzap(fmt.Errorf("could not offset tokens: %v", err))
				}
				t.Data["git-diff"] = truncatedDiff
			}

			return t.ErrorTzap(nil)
		}}
}
