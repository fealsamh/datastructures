package interned

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestInternedString(t *testing.T) {
	a := assert.New(t)

	pool := NewStringPool()

	s1 := fmt.Sprintf("abcd-%d", 1234)
	s2 := fmt.Sprintf("abcd-%d", 1234)
	is1 := pool.Get(s1)
	is2 := pool.Get(s2)

	a.True((*reflect.StringHeader)(unsafe.Pointer(&s1)).Data != (*reflect.StringHeader)(unsafe.Pointer(&s2)).Data)
	a.True((*reflect.StringHeader)(unsafe.Pointer(&is1)).Data == (*reflect.StringHeader)(unsafe.Pointer(&is2)).Data)
	a.True(is1.Eq(&is2))
}
