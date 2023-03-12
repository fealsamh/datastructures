package interned

import (
	"strings"

	"github.com/fealsamh/datastructures/unionfind"
)

// String is an interned string.
type String string

// Compare compares two interned strings.
func (s1 String) Compare(s2 String) int { return strings.Compare(string(s1), string(s2)) }

// StringPool is a pool of interned strings.
type StringPool struct {
	strings *unionfind.Structure[String]
}

// NewStringPool creates a new pool of interned strings.
func NewStringPool() *StringPool {
	return &StringPool{
		strings: unionfind.New[String](),
	}
}

// Get returns an interned strings from the pool.
func (p *StringPool) Get(s string) String {
	r, _ := p.strings.Add(String(s))
	return r.Value
}
