package reflectutil

import "reflect"

func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}
	typeZero := reflect.Zero(reflect.TypeOf(v)).Interface()
	r := reflect.DeepEqual(v, typeZero)
	return r
}
