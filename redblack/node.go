package redblack

import (
	"fmt"

	"github.com/fealsamh/datastructures/constraints"
)

type color byte

const (
	black color = iota
	red
)

type direction byte

const (
	exact direction = iota
	left
	right
)

type node[K constraints.Comparable[K], V any] struct {
	key    K
	value  V
	color  color
	parent *node[K, V]
	tree   *Tree[K, V]
	left   *node[K, V]
	right  *node[K, V]
}

func (n *node[K, V]) check() bool {
	if n.left != nil {
		if n.left.key.Compare(n.key) >= 0 {
			return false
		}
		if n.left.parent != n {
			return false
		}
		if !n.left.check() {
			return false
		}
	}
	if n.right != nil {
		if n.key.Compare(n.right.key) >= 0 {
			return false
		}
		if n.right.parent != n {
			return false
		}
		if !n.right.check() {
			return false
		}
	}
	return true
}

func (n *node[K, V]) depth() int {
	var ld, rd int
	if n.left != nil {
		ld = n.left.depth()
	}
	if n.right != nil {
		rd = n.right.depth()
	}
	if ld > rd {
		return ld + 1
	}
	return rd + 1
}

func (n *node[K, V]) size() int {
	r := 1
	if n.left != nil {
		r += n.left.size()
	}
	if n.right != nil {
		r += n.right.size()
	}
	return r
}

func (n *node[K, V]) keys() []K {
	var ks []K
	if n.left != nil {
		ks = n.left.keys()
	}
	ks = append(ks, n.key)
	if n.right != nil {
		ks = append(ks, n.right.keys()...)
	}
	return ks
}

func (n *node[K, V]) enumerate(f func(K, V)) {
	if n.left != nil {
		n.left.enumerate(f)
	}
	f(n.key, n.value)
	if n.right != nil {
		n.right.enumerate(f)
	}
}

func (n *node[K, V]) minKey() *K {
	if n.left != nil {
		return n.left.minKey()
	}
	return &n.key
}

func (n *node[K, V]) str() string {
	var s string
	if n.left != nil {
		s += "(" + n.left.str() + ") "
	}
	s += fmt.Sprintf("%v:%v", n.key, n.value)
	if n.color == black {
		s += "/B"
	} else {
		s += "/R"
	}
	if n.right != nil {
		s += " (" + n.right.str() + ")"
	}
	return s
}

func (n *node[K, V]) find(key K) (*node[K, V], direction) {
	c := key.Compare(n.key)
	switch {
	case c == 0:
		return n, exact
	case c < 0:
		if n.left == nil {
			return n, left
		}
		return n.left.find(key)
	case c > 0:
		if n.right == nil {
			return n, right
		}
		return n.right.find(key)
	}
	panic("bad red-black node")
}

func (n *node[K, V]) rotateRight() {
	p := n.parent
	pp := p.parent
	a, b, c := n.left, n.right, p.right
	if pp != nil {
		switch p.dir() {
		case left:
			pp.left = n
		case right:
			pp.right = n
		default:
			panic("bad red-black node")
		}
	} else {
		n.tree.root = n
	}
	n.parent, p.parent = pp, n
	n.left, n.right = a, p
	p.left, p.right = b, c
	if a != nil {
		a.parent = n
	}
	if b != nil {
		b.parent = p
	}
	if c != nil {
		c.parent = p
	}
}

func (n *node[K, V]) rotateLeft() {
	p := n.parent
	pp := p.parent
	a, b, c := p.left, n.left, n.right
	if pp != nil {
		switch p.dir() {
		case left:
			pp.left = n
		case right:
			pp.right = n
		default:
			panic("bad red-black node")
		}
	} else {
		n.tree.root = n
	}
	n.parent, p.parent = pp, n
	n.left, n.right = p, c
	p.left, p.right = a, b
	if c != nil {
		c.parent = n
	}
	if a != nil {
		a.parent = p
	}
	if b != nil {
		b.parent = p
	}
}

func (n *node[K, V]) rotate() {
	switch n.dir() {
	case right:
		n.rotateLeft()
	case left:
		n.rotateRight()
	}
}

func (n *node[K, V]) dir() direction {
	p := n.parent
	switch {
	case p.left == n:
		return left
	case p.right == n:
		return right
	}
	panic("bad red-black node")
}

func (n *node[K, V]) brother() *node[K, V] {
	p := n.parent
	switch {
	case p.left == n:
		return p.right
	case p.right == n:
		return p.left
	}
	panic("bad red-black node")
}

func (n *node[K, V]) ensureInvariants() {
	p := n.parent
	if p == nil {
		n.color = black
		return
	}
	if p.color == black {
		return
	}
	pp := p.parent
	if pp != nil && pp.color == black {
		u := p.brother()
		if u != nil && u.color == red {
			p.color, pp.color, u.color = black, red, black
			pp.ensureInvariants()
		} else {
			if n.dir() == p.dir() {
				p.rotate()
				p.color, pp.color = black, red
			} else {
				n.rotate()
				n.rotate()
				n.color, pp.color = black, red
			}
		}
	}
}
