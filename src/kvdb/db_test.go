package kvdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/ArborDB/arbordb/src/core"
	"github.com/ArborDB/arbordb/src/storage"
)

func TestLargeTransaction(t *testing.T) {
	store := storage.NewMemory()
	db := New(store, core.Identifier{})

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Create a transaction large enough to likely cause a stack overflow
	// with the old recursive implementation (e.g. 50,000+ operations).
	n := 10000
	for i := range n {
		key := fmt.Sprintf("k%d", i)
		val := fmt.Sprintf("v%d", i)
		if err := tx.Put(key, val); err != nil {
			t.Fatal(err)
		}
	}

	// Perform some deletions
	for i := range 100 {
		key := fmt.Sprintf("k%d", i)
		if err := tx.Delete(key); err != nil {
			t.Fatal(err)
		}
	}

	// Commit should succeed without panic/stack overflow
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	// Verify data
	tx2, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Check deleted key
	if _, err := tx2.Get("k0"); err == nil {
		t.Fatal("expected error for deleted key")
	}

	// Check existing key
	lastIdx := n - 1
	key := fmt.Sprintf("k%d", lastIdx)
	expected := fmt.Sprintf("v%d", lastIdx)
	val, err := tx2.Get(key)
	if err != nil {
		t.Fatal(err)
	}
	if val != expected {
		t.Fatalf("expected %s, got %s", expected, val)
	}
}
