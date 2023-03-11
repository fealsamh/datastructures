package logic

// Term is an n-ary term.
type Term struct {
	Symbol string
	Args   []Term
}

func (t Term) String() string {
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

// NewTerm creates a new n-ary term.
func NewTerm(symbol string, args ...string) Term {
	targs := make([]Term, len(args))
	for i, arg := range args {
		targs[i] = Term{Symbol: arg}
	}
	return Term{Symbol: symbol, Args: targs}
}
