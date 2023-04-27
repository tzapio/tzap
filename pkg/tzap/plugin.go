package tzap

import "github.com/tzapio/tzap/pkg/types"

func (t *Tzap) Plugin(pluginFunc func(tg types.TGenerator) (types.TGenerator, error)) *Tzap {
	pluginFunc(t.TG)
	return t
}
