package dshash

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"reflect"
	"slices"
)

type Kind [1]byte

var (
	KindInvalid Kind = Kind{}
	KindNil          = Kind{10}
	KindBool         = Kind{20}
	KindInt          = Kind{30}
	KindUint         = Kind{40}
	KindFloat        = Kind{50}
	KindString       = Kind{60}
	KindList         = Kind{70}
	KindListEnd      = Kind{75}
	KindMap          = Kind{80}
	KindMapEnd       = Kind{85}
	KindCycle        = Kind{90}
)

func encodeNil(ctx *Context) error {
	_, err := ctx.state.Write(KindNil[:])
	if err != nil {
		return err
	}
	return nil
}

func encodeBool(ctx *Context, value bool) error {
	_, err := ctx.state.Write(KindBool[:])
	if err != nil {
		return err
	}
	if value {
		_, err := ctx.state.Write([]byte{1})
		if err != nil {
			return err
		}
	} else {
		_, err := ctx.state.Write([]byte{0})
		if err != nil {
			return err
		}
	}
	return nil
}

func encodeInt(ctx *Context, n int64) error {
	_, err := ctx.state.Write(KindInt[:])
	if err != nil {
		return err
	}
	err = binary.Write(ctx.state, binary.LittleEndian, n)
	if err != nil {
		return err
	}
	return nil
}

func encodeUint(ctx *Context, n uint64) error {
	_, err := ctx.state.Write(KindUint[:])
	if err != nil {
		return err
	}
	err = binary.Write(ctx.state, binary.LittleEndian, n)
	if err != nil {
		return err
	}
	return nil
}

func encodeFloat(ctx *Context, n float64) error {
	_, err := ctx.state.Write(KindFloat[:])
	if err != nil {
		return err
	}
	err = binary.Write(ctx.state, binary.LittleEndian, n)
	if err != nil {
		return err
	}
	return nil
}

func encodeString(ctx *Context, value string) error {
	_, err := ctx.state.Write(KindString[:])
	if err != nil {
		return err
	}
	if err := binary.Write(ctx.state, binary.LittleEndian, int64(len(value))); err != nil {
		return err
	}
	if _, err := ctx.state.Write([]byte(value)); err != nil {
		return err
	}
	return nil
}

func encodeBytes(ctx *Context, value []byte) error {
	_, err := ctx.state.Write(KindString[:])
	if err != nil {
		return err
	}
	if err := binary.Write(ctx.state, binary.LittleEndian, int64(len(value))); err != nil {
		return err
	}
	if _, err := ctx.state.Write(value); err != nil {
		return err
	}
	return nil
}

func encodeList(ctx *Context, value reflect.Value, fn _HashFunc) error {
	_, err := ctx.state.Write(KindList[:])
	if err != nil {
		return err
	}
	for i, l := 0, value.Len(); i < l; i++ {
		if fn != nil {
			if err := fn(ctx, value.Index(i)); err != nil {
				return err
			}
		} else {
			// dynamic
			elem := value.Index(i)
			if err := getFunc(elem.Type())(ctx, elem); err != nil {
				return err
			}
		}
	}
	_, err = ctx.state.Write(KindListEnd[:])
	if err != nil {
		return err
	}
	return nil
}

func encodeStruct(ctx *Context, value reflect.Value, infos []*_FieldInfo) error {
	_, err := ctx.state.Write(KindMap[:])
	if err != nil {
		return err
	}
	for _, info := range infos {
		fieldValue := value.FieldByIndex(info.Field.Index)
		if fieldValue.IsZero() {
			// hash will not change after adding fields with zero values
			continue
		}
		if info.Func != nil {
			if err := encodeString(ctx, info.Field.Name); err != nil {
				return err
			}
			if err := info.Func(ctx, fieldValue); err != nil {
				return err
			}
		} else {
			// dynamic
			if valueIsUnsupported(fieldValue) {
				continue
			}
			if err := encodeString(ctx, info.Field.Name); err != nil {
				return err
			}
			if err := getFunc(fieldValue.Type())(ctx, fieldValue); err != nil {
				return err
			}
		}
	}
	_, err = ctx.state.Write(KindMapEnd[:])
	if err != nil {
		return err
	}
	return nil
}

type _MapEntry struct {
	Key     reflect.Value
	Value   reflect.Value
	KeyHash []byte
}

func encodeMap(ctx *Context, value reflect.Value, keyFunc _HashFunc, valueFunc _HashFunc) error {
	var entries []*_MapEntry
	iter := value.MapRange()
	for iter.Next() {
		keyState := sha256.New()
		if err := HashValue(&Context{
			state:   keyState,
			visited: ctx.visited,
		}, iter.Key()); err != nil {
			return err
		}
		entries = append(entries, &_MapEntry{
			Key:     iter.Key(),
			Value:   iter.Value(),
			KeyHash: keyState.Sum(nil),
		})
	}
	slices.SortStableFunc(entries, func(a, b *_MapEntry) int {
		return bytes.Compare(a.KeyHash, b.KeyHash)
	})

	_, err := ctx.state.Write(KindMap[:])
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if keyFunc != nil {
			if err := keyFunc(ctx, entry.Key); err != nil {
				return err
			}
		} else {
			// dynamic
			if err := getFunc(entry.Key.Type())(ctx, entry.Key); err != nil {
				return err
			}
		}
		if valueFunc != nil {
			if err := valueFunc(ctx, entry.Value); err != nil {
				return err
			}
		} else {
			// dynamic
			if err := getFunc(entry.Value.Type())(ctx, entry.Value); err != nil {
				return err
			}
		}
	}
	_, err = ctx.state.Write(KindMapEnd[:])
	if err != nil {
		return err
	}
	return nil
}

func encodeCycle(ctx *Context) error {
	_, err := ctx.state.Write(KindCycle[:])
	if err != nil {
		return err
	}
	return nil
}
