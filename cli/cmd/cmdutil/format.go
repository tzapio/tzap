package cmdutil

import (
	"fmt"

	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
)

func FormatVectorToClickable(embedding *actionpb.Embedding) string {
	return fmt.Sprintf("%s:%d", embedding.File, embedding.LineStart)
}
