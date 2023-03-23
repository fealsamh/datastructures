package unionfind

import (
	"fmt"

	"github.com/fealsamh/datastructures/constraints"
	"github.com/fealsamh/datastructures/redblack"
	"github.com/fealsamh/datastructures/sahuaro"
)

// Structure is a union-find structure.
type Structure[T constraints.Comparable[T]] struct {
	values *redblack.Tree[T, sahuaro.Tree[T]]
}

// New creates a new union-find structure.
func New[T constraints.Comparable[T]]() *Structure[T] {
	return &Structure[T]{
		values: redblack.NewTree[T, sahuaro.Tree[T]](),
	}
}

// Add adds a value to a union-find structure.
func (s *Structure[T]) Add(val T) (*sahuaro.Tree[T], bool) {
	return s.values.GetElsePut(val, func() *sahuaro.Tree[T] {
		return &sahuaro.Tree[T]{
			Value: val,
		}
	})
}

// Get retrieves an in-tree from a union-find structure.
func (s *Structure[T]) Get(val T) (*sahuaro.Tree[T], bool) {
	return s.values.Get(val)
}

// MustGet retrieves an in-tree from a union-find structure.
// It panics if the value isn't found.
func (s *Structure[T]) MustGet(val T) *sahuaro.Tree[T] {
	n, ok := s.values.Get(val)
	if !ok {
		panic(fmt.Sprintf("value '%v' not found in union-find structure", val))
	}
	return n
}
