package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"sync"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/dshash"
)

type Memory struct {
	mu    sync.RWMutex
	items map[core.Identifier]core.Expression
}

func NewMemory() *Memory {
	return &Memory{
		items: make(map[core.Identifier]core.Expression),
	}
}

var _ core.PhysicalStorage = (*Memory)(nil)

func (m *Memory) Set(ctx *core.Context, expr core.Expression) (core.Identifier, error) {
	state := sha256.New()
	if err := dshash.Hash(state, expr); err != nil {
		return core.Identifier{}, err
	}
	id := core.Identifier{
		Kind: "dshash-sha256",
		Key:  hex.EncodeToString(state.Sum(nil)),
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[id] = expr
	return id, nil
}

func (m *Memory) Get(ctx *core.Context, id core.Identifier, target any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	expr, ok := m.items[id]
	if !ok {
		return fmt.Errorf("item not found: %v", id)
	}

	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Pointer || targetVal.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	exprVal := reflect.ValueOf(expr)
	targetElem := targetVal.Elem()

	if !exprVal.Type().AssignableTo(targetElem.Type()) {
		return fmt.Errorf("cannot assign %T to %T", expr, targetElem.Interface())
	}

	targetElem.Set(exprVal)
	return nil
}
