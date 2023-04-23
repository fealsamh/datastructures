package logic

import "strings"

// Term is an n-ary term.
type Term struct {
	Symbol string
	Args   []*Term
}

func (t *Term) String() string {
	s := t.Symbol
	if len(t.Args) > 0 {
		s += "("
		for i, arg := range t.Args {
			if i > 0 {
				s += ","
			}
			s += arg.String()
		}
		s += ")"
	}
	return s
}

// Compare compares two n-ary terms using the shortlex order.
func (t *Term) Compare(t2 *Term) int {
	if c := len(t.Args) - len(t2.Args); c != 0 {
		return c
	}

	if c := strings.Compare(t.Symbol, t2.Symbol); c != 0 {
		return c
	}

	for i, arg1 := range t.Args {
		arg2 := t2.Args[i]
		if c := arg1.Compare(arg2); c != 0 {
			return c
		}
	}
	return 0
}

// NewTerm creates a new n-ary term.
func NewTerm(symbol string, args ...string) *Term {
	targs := make([]*Term, len(args))
	for i, arg := range args {
		targs[i] = &Term{Symbol: arg}
	}
	return &Term{Symbol: symbol, Args: targs}
}
