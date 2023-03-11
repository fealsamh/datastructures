package redblack

import "github.com/fealsamh/datastructures/constraints"

// Tree is a generic red-black tree.
type Tree[K constraints.Comparable[K], V any] struct {
	root *node[K, V]
}

// NewTree creates a new red-black tree.
func NewTree[K constraints.Comparable[K], V any]() *Tree[K, V] { return new(Tree[K, V]) }

// Depth returns the depth of the tree.
func (t *Tree[K, V]) Depth() int {
	if t.root == nil {
		return 0
	}
	return t.root.depth()
}

// Size returns the size of the tree.
func (t *Tree[K, V]) Size() int {
	if t.root == nil {
		return 0
	}
	return t.root.size()
}

// Keys returns the keys of the items in the tree.
func (t *Tree[K, V]) Keys() []K {
	if t.root == nil {
		return nil
	}
	return t.root.keys()
}

// Put inserts a new key-value pair into the tree or replaces the value for an existing key.
func (t *Tree[K, V]) Put(key K, value *V) (*V, bool) {
	if t.root == nil {
		t.root = &node[K, V]{key: key, value: value, color: black, tree: t}
		return nil, false
	}
	n, dir := t.root.find(key)
	switch dir {
	case exact:
		oldValue := n.value
		n.value = value
		return oldValue, true
	case left:
		l := &node[K, V]{key: key, value: value, color: red, parent: n, tree: t}
		n.left = l
		l.ensureInvariants()
	case right:
		l := &node[K, V]{key: key, value: value, color: red, parent: n, tree: t}
		n.right = l
		l.ensureInvariants()
	}
	return nil, false
}

// Get returns the value for the given key or nil if the key can't be found.
func (t *Tree[K, V]) Get(key K) (*V, bool) {
	if t.root == nil {
		return nil, false
	}
	n, dir := t.root.find(key)
	if dir == exact {
		return n.value, true
	}
	return nil, false
}

// String returns the textual representation of the tree.
func (t *Tree[K, V]) String() string {
	if t.root == nil {
		return "-"
	}
	return t.root.str()
}

// Check verifies that the keys in the tree are ordered correctly.
func (t *Tree[K, V]) Check() bool {
	if t.root == nil {
		return true
	}
	return t.root.check()
}
