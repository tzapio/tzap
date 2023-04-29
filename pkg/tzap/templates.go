package tzap

import "github.com/tzapio/tzap/pkg/types"

// ApplyTemplateP applies a given template Tzap instance to the current Tzap instance.
// Returns the applied template with its Parent set to the current Tzap instance.
func (t *Tzap) ApplyTemplateP(template *Tzap) *Tzap {
	at := t.CloneTzap(&Tzap{Name: "ApplyTemplateS"})
	Log(t, "Applying template")
	template.Parent = at
	return template
}

// ApplyTemplateFN applies a function that takes a Tzap instance and returns a modified Tzap instance.
// Returns the result of the given function applied to the current Tzap instance.
func (t *Tzap) ApplyTemplateFN(nt func(*Tzap) *Tzap) *Tzap {
	Log(t, "Applying template FN")
	return nt(t.CloneTzap(&Tzap{Name: "ApplyTemplate"}))
}
func (t *Tzap) ApplyTemplate(nt types.NamedTemplate[*Tzap]) *Tzap {
	Log(t, "Applying template FN")
	return nt.Template(t.CloneTzap(&Tzap{Name: "ApplyTemplate (" + nt.Name + ")"}))
}
