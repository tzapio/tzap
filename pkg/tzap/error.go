package tzap

import "fmt"

type ErrorTzap struct {
	Tzap *Tzap
	Err  error
}

func (t *Tzap) ErrorTzap(err error) *ErrorTzap {
	if err != nil {
		Logf(t, "Error tzap (%s)", err.Error())
		err = fmt.Errorf("error tzap %s (%d): %v", t.Name, t.Id, err)
	}

	return &ErrorTzap{
		Tzap: t,
		Err:  err,
	}
}

func (t *ErrorTzap) HandleError(cb func(*ErrorTzap) error) *Tzap {
	if t.Err != nil {
		r := cb(t)
		if r != nil {
			panic(r)
		}
	}
	return t.Tzap
}

func tzapHandlePanic(err *error) {
	if r := recover(); r != nil {
		*err = r.(error)
		println(fmt.Sprintf("Panic: %s", r.(error).Error()))
	}
}
func HandlePanic(fn func()) (err error) {
	defer tzapHandlePanic(&err)
	fn()
	return
}
