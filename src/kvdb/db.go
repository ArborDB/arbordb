package kvdb

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ArborDB/arbordb/src/collection"
	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/scalar"
)

type DB struct {
	storage core.PhysicalStorage
	mu      sync.RWMutex
	rootID  core.Identifier
}

func New(storage core.PhysicalStorage, rootID core.Identifier) *DB {
	return &DB{
		storage: storage,
		rootID:  rootID,
	}
}

func (db *DB) Begin(ctx context.Context) (*Tx, error) {
	db.mu.RLock()
	rootID := db.rootID
	db.mu.RUnlock()

	c := &core.Context{
		Context:         ctx,
		PhysicalStorage: db.storage,
	}

	var rootExpr collection.Dict[scalar.String, scalar.String]
	if rootID.Key == "" {
		rootExpr = make(collection.Map[scalar.String, scalar.String])
	} else {
		var expr core.Expression
		if err := db.storage.Get(c, rootID, &expr); err != nil {
			return nil, fmt.Errorf("load root: %w", err)
		}
		var ok bool
		rootExpr, ok = expr.(collection.Dict[scalar.String, scalar.String])
		if !ok {
			return nil, fmt.Errorf("root expression is not Dict[String, String], got %T", expr)
		}
	}

	return &Tx{
		db:         db,
		baseRootID: rootID,
		baseExpr:   rootExpr,
		mutations:  make(map[string]*string),
		ctx:        c,
	}, nil
}

type Tx struct {
	db         *DB
	baseRootID core.Identifier
	baseExpr   collection.Dict[scalar.String, scalar.String]
	mutations  map[string]*string
	ctx        *core.Context
	mu         sync.Mutex
}

func (tx *Tx) Get(key string) (string, error) {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	if val, ok := tx.mutations[key]; ok {
		if val == nil {
			return "", fmt.Errorf("key %v not found", key)
		}
		return *val, nil
	}

	val, err := tx.baseExpr.Get(tx.ctx, scalar.String(key))
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (tx *Tx) Put(key string, value string) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	tx.mutations[key] = &value
	return nil
}

func (tx *Tx) Delete(key string) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	tx.mutations[key] = nil
	return nil
}

var ErrConflict = errors.New("transaction conflict")

func (tx *Tx) Commit() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	// Materialize the base dict to a Map
	var newMap collection.Map[scalar.String, scalar.String]
	transform := collection.DictToMap[scalar.String, scalar.String]{}

	if err := transform.Apply(tx.ctx, tx.baseExpr, &newMap); err != nil {
		return fmt.Errorf("materialize: %w", err)
	}

	// Apply mutations
	for k, v := range tx.mutations {
		if v == nil {
			delete(newMap, scalar.String(k))
		} else {
			newMap[scalar.String(k)] = scalar.String(*v)
		}
	}

	// Store the new Map
	newID, err := tx.db.storage.Set(tx.ctx, newMap)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	// Update DB root
	tx.db.mu.Lock()
	defer tx.db.mu.Unlock()

	if tx.db.rootID != tx.baseRootID {
		return ErrConflict
	}

	tx.db.rootID = newID
	return nil
}
