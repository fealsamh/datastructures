package unionfind

import (
	"fmt"

	"github.com/fealsamh/datastructures/constraints"
	"github.com/fealsamh/datastructures/redblack"
)

// Structure is a union-find structure.
type Structure[T constraints.Comparable[T]] struct {
	values *redblack.Tree[T, Tree[T]]
}

// New creates a new union-find structure.
func New[T constraints.Comparable[T]]() *Structure[T] {
	return &Structure[T]{
		values: redblack.NewTree[T, Tree[T]](),
	}
}

// Add adds a value to a union-find structure.
func (s *Structure[T]) Add(val T) (*Tree[T], bool) {
	if n, ok := s.values.Get(val); ok {
		return n, true
	}
	n := &Tree[T]{
		Value: val,
	}
	s.values.Put(val, n)
	return n, false
}

// Get retrieves an in-tree from a union-find structure.
func (s *Structure[T]) Get(val T) (*Tree[T], bool) {
	return s.values.Get(val)
}

// MustGet retrieves an in-tree from a union-find structure. It panics if the value isn't found.
func (s *Structure[T]) MustGet(val T) *Tree[T] {
	n, ok := s.values.Get(val)
	if !ok {
		panic(fmt.Sprintf("value '%v' not found in union-find structure", val))
	}
	return n
}

// Tree is an in-tree.
type Tree[T any] struct {
	Value  T
	parent *Tree[T]
	rank   int
}

// Find finds the root of an in-tree.
func (n *Tree[T]) Find() *Tree[T] {
	if n.parent == nil {
		return n
	}
	r := n.parent.Find()
	n.parent = r
	return r
}

// Union merges two sets.
func (n1 *Tree[T]) Union(n2 *Tree[T]) *Tree[T] {
	x, y := n1.Find(), n2.Find()
	if x == y {
		return x
	}
	if x.rank < y.rank {
		x, y = y, x
	}
	y.parent = x
	if x.rank == y.rank {
		x.rank += 1
	}
	return x
}

func (n *Tree[T]) String() string {
	return fmt.Sprintf("%v", n.Value)
}
