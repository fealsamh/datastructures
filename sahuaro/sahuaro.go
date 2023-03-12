package sahuaro

import "fmt"

// Tree is an in-tree.
type Tree[T any] struct {
	Value  T
	parent *Tree[T]
	rank   int
}

// Find finds the root of an in-tree.
func (t *Tree[T]) Find() *Tree[T] {
	if t.parent == nil {
		return t
	}
	r := t.parent.Find()
	t.parent = r
	return r
}

// Union merges two sets.
func (t *Tree[T]) Union(t2 *Tree[T]) (*Tree[T], bool) {
	x, y := t.Find(), t2.Find()
	if x == y {
		return x, true
	}
	if x.rank < y.rank {
		x, y = y, x
	}
	y.parent = x
	if x.rank == y.rank {
		x.rank++
	}
	return x, false
}

func (t *Tree[T]) String() string {
	return fmt.Sprintf("%v", t.Value)
}
