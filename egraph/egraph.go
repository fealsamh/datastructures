package egraph

import (
	"fmt"
	"strings"

	"github.com/fealsamh/datastructures/logic"
	"github.com/fealsamh/datastructures/redblack"
	"github.com/fealsamh/datastructures/sahuaro"
	"github.com/fealsamh/datastructures/unionfind"
)

type eClassID int

func (id1 eClassID) Compare(id2 eClassID) int { return int(id1) - int(id2) }

type eNode struct {
	symbol string
	args   []*sahuaro.Tree[eClassID]
}

func (n1 *eNode) Compare(n2 *eNode) int {
	if c := len(n1.args) - len(n2.args); c != 0 {
		return c
	}
	if c := strings.Compare(n1.symbol, n2.symbol); c != 0 {
		return c
	}
	for i, arg1 := range n1.args {
		arg2 := n2.args[i]
		if c := int(arg1.Value) - int(arg2.Value); c != 0 {
			return c
		}
	}
	return 0
}

type eClass struct {
	eNodes      *redblack.Set[*eNode]
	parentNodes *redblack.Set[*eNode]
}

// Graph is an e-graph.
type Graph struct {
	maxID     int
	eClassIds *unionfind.Structure[eClassID]
	hashcons  *redblack.Tree[*eNode, *sahuaro.Tree[eClassID]]
	eClasses  *redblack.Tree[eClassID, *eClass]
}

// New creates a new e-graph.
func New() *Graph {
	return &Graph{
		eClassIds: unionfind.New[eClassID](),
		hashcons:  redblack.NewTree[*eNode, *sahuaro.Tree[eClassID]](),
		eClasses:  redblack.NewTree[eClassID, *eClass](),
	}
}

// Dump dumps the e-graph's e-classes.
func (g *Graph) Dump() {
	for _, clss := range g.Classes() {
		fmt.Println(clss)
	}
}

// Classes returns all the e-classes of the e-graph.
func (g *Graph) Classes() [][]*logic.Term {
	processed := make(map[*eClass]struct{})
	var r [][]*logic.Term
	for _, id := range g.eClasses.Keys() {
		cls, _ := g.eClasses.Get(id)
		if _, ok := processed[cls]; ok {
			continue
		}
		processed[cls] = struct{}{}
		terms := redblack.NewSet[*logic.Term]()
		for _, n := range cls.eNodes.Values() {
			terms.Insert(g.getTerm(n))
		}
		r = append(r, terms.Values())
	}
	return r
}

// Merge merges two n-ary terms.
func (g *Graph) Merge(t1, t2 *logic.Term) {
	_, clsID1, ok := g.getENode(t1, false)
	if !ok {
		panic(fmt.Sprintf("term '%s' not found in e-graph", t1))
	}
	_, clsID2, ok := g.getENode(t2, false)
	if !ok {
		panic(fmt.Sprintf("term '%s' not found in e-graph", t2))
	}
	g.merge(clsID1, clsID2)
}

func (g *Graph) merge(clsID1, clsID2 *sahuaro.Tree[eClassID]) {
	if _, equiv := clsID1.Union(clsID2); equiv {
		return
	}
	cls1, _ := g.eClasses.Get(clsID1.Value)
	cls2, _ := g.eClasses.Get(clsID2.Value)
	for _, cls := range cls2.eNodes.Values() {
		cls1.eNodes.Insert(cls)
	}
	for _, cls := range cls2.parentNodes.Values() {
		cls1.parentNodes.Insert(cls)
	}
	g.eClasses.Enumerate(func(id eClassID, cls *eClass) bool {
		if cls == cls2 {
			g.eClasses.Put(id, cls1)
		}
		return true
	})
	parentNodes := cls1.parentNodes.Values()
	// preserving the congruence invariant
	for i, n1 := range parentNodes {
		for j := i + 1; j < len(parentNodes); j++ {
			n2 := parentNodes[j]
			if n1.symbol == n2.symbol && len(n1.args) == len(n2.args) {
				for k, arg1 := range n1.args {
					arg2 := n2.args[k]
					if arg1.Find() != arg2.Find() {
						return
					}
				}
				id1, _ := g.hashcons.Get(n1)
				id2, _ := g.hashcons.Get(n2)
				g.merge(id1, id2)
			}
		}
	}
}

// Get retrieves the representative of an n-ary term from the e-graph.
func (g *Graph) Get(t *logic.Term) (*logic.Term, bool) {
	n, _, ok := g.getENode(t, false)
	if !ok {
		return nil, false
	}
	clsID, _ := g.hashcons.Get(n)
	cls, _ := g.eClasses.Get(clsID.Value)
	n = *cls.eNodes.MinKey()
	return g.getTerm(n), true
}

func (g *Graph) getTerm(n *eNode) *logic.Term {
	args := make([]*logic.Term, len(n.args))
	for i, arg := range n.args {
		cls, ok := g.eClasses.Get(arg.Value)
		if !ok {
			panic("e-class must exist at this point")
		}
		args[i] = g.getTerm(*cls.eNodes.MinKey())
	}
	return &logic.Term{Symbol: n.symbol, Args: args}
}

// Add adds an n-ary term to the e-graph.
func (g *Graph) Add(t *logic.Term) bool {
	_, _, ok := g.getENode(t, true)
	return ok
}

func (g *Graph) getEClassID(n *eNode, create bool) (*sahuaro.Tree[eClassID], bool) {
	clsID, ok := g.hashcons.Get(n)
	if !ok {
		if create {
			g.maxID++
			clsID := eClassID(g.maxID)
			t, _ := g.eClassIds.Add(clsID)
			g.hashcons.Put(n, t)
			cls := &eClass{
				eNodes:      redblack.NewSet[*eNode](),
				parentNodes: redblack.NewSet[*eNode](),
			}
			cls.eNodes.Insert(n)
			g.eClasses.Put(clsID, cls)
			return t, false
		}
		return nil, false
	}
	return clsID.Find(), true
}

func (g *Graph) getENode(t *logic.Term, create bool) (*eNode, *sahuaro.Tree[eClassID], bool) {
	args := make([]*sahuaro.Tree[eClassID], len(t.Args))
	for i, arg := range t.Args {
		n, clsID, ok := g.getENode(arg, create)
		if !ok && !create {
			return n, nil, false
		}
		args[i] = clsID
	}
	n := &eNode{symbol: t.Symbol, args: args}
	clsID, ok := g.getEClassID(n, create)
	if !ok && create {
		for _, arg := range n.args {
			cls, _ := g.eClasses.Get(arg.Value)
			cls.parentNodes.Insert(n)
		}
	}
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
		if !g.IsCanonicalEClassID(arg.Value) {
			return false
		}
	}
	return true
}

// CheckEClassMap checks whether the e-class map is valid.
func (g *Graph) CheckEClassMap() bool {
	type pair struct {
		canonicalEClassID *sahuaro.Tree[eClassID]
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
