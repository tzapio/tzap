package tzap

// ChangeFilepath updates the filepath metadata in the Tzap data
func (t *Tzap) ChangeFilepath(filepath string) *Tzap {
	t.Data["filepath"] = filepath
	return t
}
