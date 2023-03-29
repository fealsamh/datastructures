package redblack

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
)

type pair[T, U any] struct {
	fst T
	snd U
}

func newPair[T, U any](fst T, snd U) pair[T, U] {
	return pair[T, U]{
		fst: fst,
		snd: snd,
	}
}

func generateTestData() []pair[string, int] {
	var r []pair[string, int]
	for i := 1; i <= 1_000; i++ {
		n := rand.Intn(1_000_000)
		r = append(r, newPair(fmt.Sprintf("k%d", n), n))
	}
	return r
}

var gR interface{}

func BenchmarkBuiltinMap(b *testing.B) {
	pairs := generateTestData()
	var lR interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m := make(map[string]int)
		for _, p := range pairs {
			m[p.fst] = p.snd
		}
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		lR = newPair(m, keys)
	}
	gR = lR
}

func BenchmarkRedblackTree(b *testing.B) {
	pairs := generateTestData()
	var lR interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		t := NewTree[compString, int]()
		for _, p := range pairs {
			t.Put(compString(p.fst), &p.snd)
		}
		keys := make([]string, 0, t.Size())
		t.Enumerate(func(k compString, _ *int) {
			keys = append(keys, string(k))
		})
		lR = newPair(t, keys)
	}
	gR = lR
}

type compString string

func (s1 compString) Compare(s2 compString) int {
	return strings.Compare(string(s1), string(s2))
}
