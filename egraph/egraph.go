package egraph

import (
	"strings"

	"github.com/fealsamh/datastructures/logic"
	"github.com/fealsamh/datastructures/redblack"
	"github.com/fealsamh/datastructures/unionfind"
)

type eClassID int

func (id1 eClassID) Compare(id2 eClassID) int { return int(id1) - int(id2) }

type eNode struct {
	symbol string
	args   []eClassID
}

func (n1 eNode) Compare(n2 eNode) int {
	c := len(n1.args) - len(n2.args)
	if c != 0 {
		return c
	}
	c = strings.Compare(n1.symbol, n2.symbol)
	if c != 0 {
		return c
	}
	for i, arg1 := range n1.args {
		arg2 := n2.args[i]
		c := int(arg1) - int(arg2)
		if c != 0 {
			return c
		}
	}
	return 0
}

type eClass redblack.Set[eNode]

// Graph is an e-graph.
type Graph struct {
	maxID     int
	eClassIds *unionfind.Structure[eClassID]
	hashcons  *redblack.Tree[eNode, eClassID]
	eClasses  *redblack.Tree[eClassID, eClass]
}

// New create a new e-graph.
func New() *Graph {
	return &Graph{
		eClassIds: unionfind.New[eClassID](),
		hashcons:  redblack.NewTree[eNode, eClassID](),
		eClasses:  redblack.NewTree[eClassID, eClass](),
	}
}

// Get retrieves the representative of an n-ary term from the e-graph.
func (g *Graph) Get(t logic.Term) (*logic.Term, bool) {
	n, _, ok := g.getENode(&t, false)
	if !ok {
		return nil, false
	}
	return g.getTerm(&n), true
}

func (g *Graph) getTerm(n *eNode) *logic.Term {
	args := make([]logic.Term, len(n.args))
	for i, arg := range n.args {
		cls, _ := g.eClasses.Get(arg)
		args[i] = *g.getTerm((*redblack.Set[eNode])(cls).MinKey())
	}
	return &logic.Term{Symbol: n.symbol, Args: args}
}

// Add adds an n-ary term to the e-graph.
func (g *Graph) Add(t logic.Term) bool {
	_, _, ok := g.getENode(&t, true)
	return ok
}

func (g *Graph) getEClassID(n eNode, create bool) (eClassID, bool) {
	clsID, ok := g.hashcons.Get(n)
	if !ok {
		if create {
			g.maxID++
			clsID := eClassID(g.maxID)
			g.eClassIds.Add(clsID)
			g.hashcons.Put(n, &clsID)
			cls := redblack.NewSet[eNode]()
			cls.Insert(n)
			g.eClasses.Put(clsID, (*eClass)(cls))
			return clsID, false
		}
		return 0, false
	}
	return g.eClassIds.MustGet(*clsID).Find().Value, true
}

func (g *Graph) getENode(t *logic.Term, create bool) (eNode, eClassID, bool) {
	args := make([]eClassID, len(t.Args))
	for i, arg := range t.Args {
		n, clsID, ok := g.getENode(&arg, create)
		if !ok && !create {
			return n, 0, false
		}
		args[i] = clsID
	}
	n := eNode{symbol: t.Symbol, args: args}
	clsID, ok := g.getEClassID(n, create)
	return n, clsID, ok
}

// IsCanonicalEClassID determines whether `id` is canonical.
func (g *Graph) IsCanonicalEClassID(id eClassID) bool {
	t := g.eClassIds.MustGet(id)
	return t.Find() == t
}

// IsCanonicalENode determines whether `n` is canonical.
func (g *Graph) IsCanonicalENode(n eNode) bool {
	for _, arg := range n.args {
		if !g.IsCanonicalEClassID(arg) {
			return false
		}
	}
	return true
}

// CheckEClassMap checks whether the e-class map is valid.
func (g *Graph) CheckEClassMap() bool {
	type pair struct {
		canonicalEClassID *unionfind.Tree[eClassID]
		cls               *eClass
	}
	var pairs []pair
	for _, id := range g.eClasses.Keys() {
		cls, _ := g.eClasses.Get(id)
		pairs = append(pairs, pair{
			canonicalEClassID: g.eClassIds.MustGet(id).Find(),
			cls:               cls,
		})
	}
	for i, p1 := range pairs {
		for j := i + 1; j < len(pairs); j++ {
			p2 := pairs[j]
			if (p1.canonicalEClassID == p2.canonicalEClassID) != (p1.cls == p2.cls) {
				return false
			}
		}
	}
	return true
}
