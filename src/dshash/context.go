package dshash

import (
	"hash"
	"unsafe"
)

type Context struct {
	state   hash.Hash
	visited map[unsafe.Pointer]struct{}
}
