package cmdutil

import (
	"fmt"

	"github.com/tzapio/tzap/pkg/types"
)

func FormatVectorToClickable(v types.Vector) string {
	return fmt.Sprintf("%s:%d", v.Metadata.Filename, v.Metadata.LineStart)
}
