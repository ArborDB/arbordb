package dshash

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"sync"
)

type _HashFunc = func(ctx *Context, value reflect.Value) error

var funcs sync.Map // reflect.Type -> _HashFunc

type _FieldInfo struct {
	Field reflect.StructField
	Func  _HashFunc
	Once  sync.Once
}

func makeFunc(t reflect.Type) _HashFunc {
	if isUnsupported(t) {
		return func(ctx *Context, value reflect.Value) error {
			return encodeNil(ctx)
		}
	}

	switch t.Kind() {

	case reflect.Pointer:
		return func(ctx *Context, value reflect.Value) error {
			if value.IsNil() {
				return encodeNil(ctx)
			}
			if _, ok := ctx.visited[value.UnsafePointer()]; ok {
				return encodeCycle(ctx)
			}
			ctx.visited[value.UnsafePointer()] = struct{}{}
			elemFunc := getFunc(t.Elem())
			return elemFunc(ctx, value.Elem())
		}

	case reflect.Interface:
		return func(ctx *Context, value reflect.Value) error {
			if value.IsNil() {
				return encodeNil(ctx)
			}
			elem := value.Elem()
			return getFunc(elem.Type())(ctx, elem)
		}

	case reflect.Bool:
		return func(ctx *Context, value reflect.Value) error {
			return encodeBool(ctx, value.Bool())
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(ctx *Context, value reflect.Value) error {
			return encodeInt(ctx, value.Int())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(ctx *Context, value reflect.Value) error {
			return encodeUint(ctx, value.Uint())
		}

	case reflect.Float32, reflect.Float64:
		return func(ctx *Context, value reflect.Value) error {
			return encodeFloat(ctx, value.Float())
		}

	case reflect.String:
		return func(ctx *Context, value reflect.Value) error {
			return encodeString(ctx, value.String())
		}

	case reflect.Array, reflect.Slice:
		elemType := t.Elem()
		if elemType == byteType {
			// []byte or [...]byte
			return func(ctx *Context, value reflect.Value) error {
				if value.Kind() == reflect.Array && !value.CanAddr() {
					// reflect.Value.Bytes() panics if the array is not addressable.
					// We must copy it to a slice to access the bytes safely.
					tmp := make([]byte, value.Len())
					reflect.Copy(reflect.ValueOf(tmp), value)
					return encodeBytes(ctx, tmp)
				}
				return encodeBytes(ctx, value.Bytes())
			}
		}
		var fn _HashFunc
		var once sync.Once
		return func(ctx *Context, value reflect.Value) error {
			once.Do(func() {
				if elemType.Kind() != reflect.Interface {
					fn = getFunc(elemType)
				}
			})
			return encodeList(ctx, value, fn)
		}

	case reflect.Struct:
		var infos []*_FieldInfo
		for i := range t.NumField() {
			field := t.Field(i)
			if isUnsupported(field.Type) {
				continue
			}
			infos = append(infos, &_FieldInfo{
				Field: field,
			})
		}
		slices.SortFunc(infos, func(a, b *_FieldInfo) int {
			return cmp.Compare(a.Field.Name, b.Field.Name)
		})
		return func(ctx *Context, value reflect.Value) error {
			return encodeStruct(ctx, value, infos)
		}

	case reflect.Map:
		var keyFunc, valueFunc _HashFunc
		var once sync.Once
		return func(ctx *Context, value reflect.Value) error {
			once.Do(func() {
				if t.Key().Kind() != reflect.Interface {
					keyFunc = getFunc(t.Key())
				}
				if t.Elem().Kind() != reflect.Interface {
					valueFunc = getFunc(t.Elem())
				}
			})
			return encodeMap(ctx, value, keyFunc, valueFunc)
		}

	}

	panic(fmt.Errorf("unknown type: %v", t))
}

func getFunc(t reflect.Type) _HashFunc {
	v, ok := funcs.Load(t)
	if ok {
		return v.(_HashFunc)
	}
	fn := makeFunc(t)
	v, _ = funcs.LoadOrStore(t, fn)
	return v.(_HashFunc)
}
