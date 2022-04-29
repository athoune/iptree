package tree

import (
	"net"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestTree2(t *testing.T) {
	tree := NewTrunk(2)
	_, nm, err := net.ParseCIDR("192.168.1.0/24")
	assert.NoError(t, err)
	tree.Append(nm, "Hello")
	spew.Dump(tree)
	_, ok := tree.Get(net.ParseIP("192.168.1.42"))
	assert.True(t, ok)
	_, ok = tree.Get(net.ParseIP("192.168.2.42"))
	assert.False(t, ok)
}

func BenchmarkContains(b *testing.B) {
	_, nm, err := net.ParseCIDR("192.168.1.0/24")
	assert.NoError(b, err)
	a := net.ParseIP("192.168.1.42")
	for i := 0; i < b.N; i++ {
		nm.Contains(a)
	}
}
