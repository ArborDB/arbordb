package dshash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
)

func ptrTo[T any](v T) *T {
	return &v
}

func TestHash(t *testing.T) {
	type Case struct {
		value any
		sum   string
	}

	for _, _case := range []Case{

		// bool
		{
			value: true,
			sum:   "29fc79494a0048a5dce25209b1d23ed2e079a255237dfab548acec4e8916cdac",
		},
		{
			value: false,
			sum:   "5aeca385d8b781825b07bbec7c858b7170426c88088935850bc13dd6402368a5",
		},

		// pointer and interface
		{
			value: ptrTo(false),
			sum:   "5aeca385d8b781825b07bbec7c858b7170426c88088935850bc13dd6402368a5",
		},
		{
			value: ptrTo((any)(false)),
			sum:   "5aeca385d8b781825b07bbec7c858b7170426c88088935850bc13dd6402368a5",
		},
		{
			value: ptrTo(ptrTo((any)(false))),
			sum:   "5aeca385d8b781825b07bbec7c858b7170426c88088935850bc13dd6402368a5",
		},
		{
			value: ptrTo((any)(ptrTo(ptrTo((any)(false))))),
			sum:   "5aeca385d8b781825b07bbec7c858b7170426c88088935850bc13dd6402368a5",
		},

		// integer
		{
			value: 42,
			sum:   "719cbd638a8b35f5be570d08c1c402f7e1277c5ef363c090cff65e6b7468ebdb",
		},
		{
			value: int8(42),
			sum:   "719cbd638a8b35f5be570d08c1c402f7e1277c5ef363c090cff65e6b7468ebdb",
		},

		// unsigned integer
		{
			value: uint64(42),
			sum:   "140ccf58356b7f90f7e93fdd5e1eb4407b54dd0d8cbb0395f52c29b1d2ae7fba",
		},
		{
			value: uint32(42),
			sum:   "140ccf58356b7f90f7e93fdd5e1eb4407b54dd0d8cbb0395f52c29b1d2ae7fba",
		},

		// float
		{
			value: float32(42),
			sum:   "4475aa3302321ad276fbfeb616c9e31b668b9a2d1e0b98abaa19e2337058d300",
		},
		{
			value: float64(42),
			sum:   "4475aa3302321ad276fbfeb616c9e31b668b9a2d1e0b98abaa19e2337058d300",
		},

		// string
		{
			value: "foo",
			sum:   "57e5773169334a5b22bd1e0bb1d7ff09522cd2ea978a4e0c42fea1ceef605fad",
		},

		// array and slice
		{
			value: []int{42, 1},
			sum:   "193e34c2e4dc72b1da2864a71bc694b7ba82fa5e4f4dd24f4784b6a16f1a86c6",
		},
		{
			value: [2]int{42, 1},
			sum:   "193e34c2e4dc72b1da2864a71bc694b7ba82fa5e4f4dd24f4784b6a16f1a86c6",
		},
		{
			value: [2]int32{42, 1},
			sum:   "193e34c2e4dc72b1da2864a71bc694b7ba82fa5e4f4dd24f4784b6a16f1a86c6",
		},
		{
			value: [2]any{42, 1},
			sum:   "193e34c2e4dc72b1da2864a71bc694b7ba82fa5e4f4dd24f4784b6a16f1a86c6",
		},
		{
			value: []any{int64(42), int16(1)},
			sum:   "193e34c2e4dc72b1da2864a71bc694b7ba82fa5e4f4dd24f4784b6a16f1a86c6",
		},

		// struct
		{
			value: struct {
				A int
				B int
			}{
				A: 42,
				B: 1,
			},
			sum: "27fb15b6e91ee74e320afc0361ae9c929fdf145734eb231c47a15547dfc3d65d",
		},
		{
			value: struct {
				A int
				B int
				C string
			}{
				A: 42,
				B: 1,
				C: "", // zero C
			},
			sum: "27fb15b6e91ee74e320afc0361ae9c929fdf145734eb231c47a15547dfc3d65d",
		},
		{
			value: struct {
				A int
				B int
				C string
			}{
				A: 42,
				B: 1,
				C: "foo", // non-zero C
			},
			sum: "88a0faef79f4bd8dc116d709ab1990f457f04bb7586e80b98c35ea0386c37950",
		},
		{
			value: struct {
				// different field order
				B int
				C string
				A int
			}{
				A: 42,
				B: 1,
				C: "foo",
			},
			sum: "88a0faef79f4bd8dc116d709ab1990f457f04bb7586e80b98c35ea0386c37950",
		},
		{
			value: struct {
				B int
				C string
				A int
				D func() // unsupported type
			}{
				A: 42,
				B: 1,
				// non-zero C
				C: "foo",
				D: func() {},
			},
			sum: "88a0faef79f4bd8dc116d709ab1990f457f04bb7586e80b98c35ea0386c37950",
		},
		{
			value: struct {
				B int
				C string
				A int
				D any // dynamic unsupported type
			}{
				A: 42,
				B: 1,
				C: "foo",
				D: func() {},
			},
			sum: "88a0faef79f4bd8dc116d709ab1990f457f04bb7586e80b98c35ea0386c37950",
		},

		// map
		{
			value: map[string]any{
				"A": 42,
				"B": 1,
				"C": "foo",
			},
			sum: "88a0faef79f4bd8dc116d709ab1990f457f04bb7586e80b98c35ea0386c37950",
		},
		{
			value: map[string]any{
				"A": 42,
				"B": 1,
				"C": "foo",
				"D": nil, // nil value
			},
			sum: "1e932111bc5736f6d7b496c5132bda1e4c3e0dd693d59cfe05c25d1687be7ee1",
		},

		// unsupported types
		{
			value: func() {},
			sum:   "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
		},
		{
			value: make(chan bool),
			sum:   "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
		},
		{
			value: ptrTo(make(chan bool)),
			sum:   "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
		},
		{
			value: ptrTo((any)(make(chan bool))),
			sum:   "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
		},
		{
			value: []any{
				func() {},
				func() {},
			},
			sum: "0f56fbbf7f975a7fc56f4f7e364592ee6c6f87877f7f090b3d3d1f954b7f1c21",
		},
		{
			value: []any{
				ptrTo(make(chan bool)),
				func() {},
			},
			sum: "0f56fbbf7f975a7fc56f4f7e364592ee6c6f87877f7f090b3d3d1f954b7f1c21",
		},
	} {

		state := sha256.New()
		if err := Hash(state, _case.value); err != nil {
			t.Fatal(err)
		}
		sum := hex.EncodeToString(state.Sum(nil))
		if sum != _case.sum {
			t.Fatalf("mismatch: %+v, got %v", _case, sum)
		}
	}
}

func TestHashNil(t *testing.T) {
	if err := Hash(sha256.New(), nil); err != nil {
		t.Fatal(err)
	}
	if err := Hash(sha256.New(), (map[int]int)(nil)); err != nil {
		t.Fatal(err)
	}
	if err := Hash(sha256.New(), ([]int)(nil)); err != nil {
		t.Fatal(err)
	}
}

func TestHashCycle(t *testing.T) {
	type P *P
	var p P
	p = &p
	h := sha256.New()
	if err := Hash(h, p); err != nil {
		t.Fatal()
	}
	sum1 := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	if err := Hash(h, &p); err != nil {
		t.Fatal()
	}
	sum2 := fmt.Sprintf("%x", h.Sum(nil))
	if sum1 != sum2 {
		t.Fatal()
	}
}
