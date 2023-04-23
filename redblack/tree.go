package redblack

import "github.com/fealsamh/datastructures/constraints"

// Tree is a generic red-black tree.
type Tree[K constraints.Comparable[K], V any] struct {
	root *node[K, V]
}

// NewTree creates a new red-black tree.
func NewTree[K constraints.Comparable[K], V any]() *Tree[K, V] { return new(Tree[K, V]) }

// NewTreeFromMap creates a new red-black tree from an in-built map.
func NewTreeFromMap[K interface {
	constraints.Comparable[K]
	comparable
}, V any](m map[K]V) *Tree[K, V] {
	tr := NewTree[K, V]()
	for k, v := range m {
		tr.Put(k, v)
	}
	return tr
}

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

// Enumerate enumerates all the items in the tree.
func (t *Tree[K, V]) Enumerate(f func(K, V) bool) bool {
	if t.root == nil {
		return true
	}
	return t.root.enumerate(f)
}

// MinKey returns the minimum key in the tree or nil if the tree is empty.
func (t *Tree[K, V]) MinKey() *K {
	if t.root == nil {
		return nil
	}
	return t.root.minKey()
}

// Put inserts a new key-value pair into the tree or replaces the value for an existing key.
func (t *Tree[K, V]) Put(key K, value V) (oldValue V, updated bool) {
	if t.root == nil {
		t.root = &node[K, V]{key: key, value: value, color: black, tree: t}
		return
	}
	n, dir := t.root.find(key)
	switch dir {
	case exact:
		oldValue = n.value
		n.value = value
		updated = true
		return
	case left:
		l := &node[K, V]{key: key, value: value, color: red, parent: n, tree: t}
		n.left = l
		l.ensureInvariants()
	case right:
		l := &node[K, V]{key: key, value: value, color: red, parent: n, tree: t}
		n.right = l
		l.ensureInvariants()
	}
	return
}

// GetElsePut returns the value for the given key or inserts a new key-value pair into the tree.
func (t *Tree[K, V]) GetElsePut(key K, fnValue func() V) (retValue V, updated bool) {
	if t.root == nil {
		retValue = fnValue()
		t.root = &node[K, V]{key: key, value: retValue, color: black, tree: t}
		return
	}
	n, dir := t.root.find(key)
	switch dir {
	case exact:
		return n.value, true
	case left:
		retValue = fnValue()
		l := &node[K, V]{key: key, value: retValue, color: red, parent: n, tree: t}
		n.left = l
		l.ensureInvariants()
		return
	case right:
		retValue = fnValue()
		l := &node[K, V]{key: key, value: retValue, color: red, parent: n, tree: t}
		n.right = l
		l.ensureInvariants()
		return
	}
	return
}

// Get returns the value for the given key or nil if the key can't be found.
func (t *Tree[K, V]) Get(key K) (retValue V, found bool) {
	if t.root == nil {
		return
	}
	n, dir := t.root.find(key)
	if dir == exact {
		return n.value, true
	}
	return
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
