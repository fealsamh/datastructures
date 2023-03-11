package redblack

import "github.com/fealsamh/datastructures/constraints"

// Set is a generic red-black tree.
type Set[K constraints.Comparable[K]] Tree[K, struct{}]

// NewSet creates a new red-black set.
func NewSet[K constraints.Comparable[K]]() *Set[K] { return (*Set[K])(NewTree[K, struct{}]()) }

// Depth returns the depth of the set.
func (s *Set[K]) Depth() int {
	return (*Tree[K, struct{}])(s).Depth()
}

// Size returns the size of the tree.
func (s *Set[K]) Size() int {
	return (*Tree[K, struct{}])(s).Size()
}

// Keys returns the elements of the set.
func (s *Set[K]) Keys() []K {
	return (*Tree[K, struct{}])(s).Keys()
}

// MinKey returns the minimum element from the set or nil if the set is empty.
func (s *Set[K]) MinKey() *K {
	return (*Tree[K, struct{}])(s).MinKey()
}

// Insert inserts a new element into the set.
func (s *Set[K]) Insert(key K) bool {
	_, ok := (*Tree[K, struct{}])(s).Put(key, nil)
	return ok
}

// Contains returns true if the element is found in the set.
func (s *Set[K]) Contains(key K) bool {
	_, ok := (*Tree[K, struct{}])(s).Get(key)
	return ok
}
