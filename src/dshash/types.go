package dshash

import "reflect"

func isUnsupported(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Invalid,
		reflect.Uintptr, reflect.UnsafePointer,
		reflect.Complex64, reflect.Complex128,
		reflect.Chan, reflect.Func:
		return true
	}
	return false
}

var byteType = reflect.TypeFor[byte]()

func valueIsUnsupported(v reflect.Value) bool {
	for (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) &&
		!v.IsNil() {
		v = v.Elem()
	}
	return isUnsupported(v.Type())
}
