package tzap

// ApplyTemplate applies a given template Tzap instance to the current Tzap instance.
// Returns the applied template with its Parent set to the current Tzap instance.
func (t *Tzap) ApplyTemplate(template *Tzap) *Tzap {
	Log(t, "Applying template")
	template.Parent = t
	return template
}

// ApplyTemplateFN applies a function that takes a Tzap instance and returns a modified Tzap instance.
// Returns the result of the given function applied to the current Tzap instance.
func (t *Tzap) ApplyTemplateFN(nt func(*Tzap) *Tzap) *Tzap {
	Log(t, "Applying template FN")
	return nt(t)
}
