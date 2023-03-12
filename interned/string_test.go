package interned

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestInternedString(t *testing.T) {
	a := assert.New(t)

	pool := NewStringPool()

	s1 := pool.Get("abcd")
	s2 := pool.Get("abcd")

	a.True((*reflect.StringHeader)(unsafe.Pointer(&s1)).Data == (*reflect.StringHeader)(unsafe.Pointer(&s2)).Data)
}
