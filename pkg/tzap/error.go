package tzap

import (
	"fmt"
	"os"
	"runtime"
)

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
		stack := make([]byte, 4096)
		length := runtime.Stack(stack, true)

		// Print the error message
		println(r.(error).Error())
		// Print the stack trace
		fmt.Fprintf(os.Stderr, "%s\n", stack[:length])
	}
}
func HandlePanic(fn func()) (err error) {
	defer tzapHandlePanic(&err)
	fn()
	return
}
