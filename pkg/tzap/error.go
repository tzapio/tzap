package tzap

type ErrorTzap struct {
	Tzap *Tzap
	Err  error
}

func (t *Tzap) ErrorTzap(err error) *ErrorTzap {
	if err != nil {
		Logf(t, "Error tzap (%s)", err.Error())
	}

	return &ErrorTzap{
		Tzap: t,
		Err:  err,
	}
}

func (t *ErrorTzap) HandleError(cb func(*ErrorTzap) *Tzap) *Tzap {
	if t.Err != nil {
		r := cb(t)
		if r == nil {
			panic(t.Err)
		}
	}
	return t.Tzap
}

func tzapHandlePanic(err *error) {
	if r := recover(); r != nil {
		println("Hello world!!!", r.(error).Error())
		*err = r.(error)
	}
}
func HandlePanic(fn func()) (err error) {
	defer tzapHandlePanic(&err)
	fn()
	return
}
