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

	for i, _case := range []Case{

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
		}, // 5
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
		}, // 10

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
			sum:   "c2b6962cc1364738a0a5e4ab23a856100ed4297eee43bee139012142aab26276",
		},

		// array and slice
		{
			value: []int{42, 1},
			sum:   "193e34c2e4dc72b1da2864a71bc694b7ba82fa5e4f4dd24f4784b6a16f1a86c6",
		},
		{
			value: [2]int{42, 1},
			sum:   "193e34c2e4dc72b1da2864a71bc694b7ba82fa5e4f4dd24f4784b6a16f1a86c6",
		}, // 15
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
			sum: "66784164dbe8d890564ef2f1fb0ec5e5ab98daf9deb5ccee1e455fc1ef4d83a9",
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
			sum: "66784164dbe8d890564ef2f1fb0ec5e5ab98daf9deb5ccee1e455fc1ef4d83a9",
		}, // 20
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
			sum: "61eb6e17bb030eca37c3c3c980c64bb8bdd17550abcde1c66668d0cd878a2ba1",
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
			sum: "61eb6e17bb030eca37c3c3c980c64bb8bdd17550abcde1c66668d0cd878a2ba1",
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
			sum: "61eb6e17bb030eca37c3c3c980c64bb8bdd17550abcde1c66668d0cd878a2ba1",
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
			sum: "61eb6e17bb030eca37c3c3c980c64bb8bdd17550abcde1c66668d0cd878a2ba1",
		},

		// map
		{
			value: map[string]any{
				"A": 42,
				"B": 1,
				"C": "foo",
			},
			sum: "3f13ec63968c60912ebbe37f28ebb024d340f1b13b5dd14d9e68f87b5bd20e65",
		}, // 25
		{
			value: map[string]any{
				"A": 42,
				"B": 1,
				"C": "foo",
				"D": nil, // nil value
			},
			sum: "6c55ea82aca985825d1d14fd6760643b1eab819383a15752c8ec6f59f37b1d1c",
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
			t.Fatalf("%d, mismatch: %+v, got %v", i+1, _case, sum)
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

func TestHashCollision_PointerReuse(t *testing.T) {
	i := 42
	p := &i
	v1 := []*int{p}
	v2 := []*int{p, p}

	h1 := sha256.New()
	Hash(h1, v1)
	sum1 := hex.EncodeToString(h1.Sum(nil))

	h2 := sha256.New()
	Hash(h2, v2)
	sum2 := hex.EncodeToString(h2.Sum(nil))

	if sum1 == sum2 {
		t.Fatal("collision detected for pointer reuse")
	}
}

func TestHashCollision_StringInjection(t *testing.T) {
	v1 := []string{"\x3c"}
	v2 := []string{"", ""}

	h1 := sha256.New()
	Hash(h1, v1)
	sum1 := hex.EncodeToString(h1.Sum(nil))

	h2 := sha256.New()
	Hash(h2, v2)
	sum2 := hex.EncodeToString(h2.Sum(nil))

	if sum1 == sum2 {
		t.Fatal("collision detected for string injection")
	}
}

func TestHashNonAddressableArray(t *testing.T) {
	// A struct with an array field
	type S struct {
		Arr [4]byte
	}

	// Hash operates on interface{}, so passing S{} passes the struct by value.
	// Inside reflection, the struct Value is not addressable, and accessing
	// the Arr field yields a non-addressable array Value.
	// This previously triggered a panic in reflect.Value.Bytes().
	if err := Hash(sha256.New(), S{Arr: [4]byte{1, 2, 3, 4}}); err != nil {
		t.Fatal(err)
	}
}
