package cmdutil

import "github.com/tzapio/tzap/pkg/types"

func FormatVectorToClickable(v types.Vector) string {
	return v.Metadata["filename"] + ":" + v.Metadata["lineStart"]
}
