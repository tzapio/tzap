package cmdutil

import (
	"context"

	"github.com/tzapio/tzap/pkg/tzap"
)

var tzapContextKey = struct{}{}

func SetTzapInContext(ctx context.Context, t *tzap.Tzap) context.Context {
	return context.WithValue(ctx, tzapContextKey, t)
}
func GetTzapFromContext(ctx context.Context) *tzap.Tzap {
	return ctx.Value(tzapContextKey).(*tzap.Tzap)
}
