package dshash

import (
	"hash"
	"reflect"
	"unsafe"
)

// deterministic structural hashing

func Hash(state hash.Hash, value any) error {
	ctx := &Context{
		state:   state,
		visited: make(map[unsafe.Pointer]struct{}),
	}
	return HashValue(ctx, reflect.ValueOf(value))
}

func HashValue(ctx *Context, value reflect.Value) error {
	if !value.IsValid() {
		return encodeNil(ctx)
	}
	return getFunc(value.Type())(ctx, value)
}
