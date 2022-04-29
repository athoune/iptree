package tree

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	lru "github.com/hashicorp/golang-lru"
)

/*
Trunk is the start of a tree storage
*/
type Trunk interface {
	// Append an IP network, and a payload
	Append(nm *net.IPNet, data interface{})
	// Get an IP, return a payload, and a boolean, true if the IP is in a registered network
	Get(ip net.IP) (interface{}, bool)
	// Size is the number of networks
	Size() int
	// Dump the tree, for debugging purpose
	Dump(w io.Writer)
}

/*
SimpleTrunk is a naïve Trunk, without any cache
*/
type SimpleTrunk struct {
	*Node
	size             int
	numberOfFullList int
}

/*
NewTrunk returns a new Trunk
*/
func NewTrunk(numberOfFullList int) *SimpleTrunk {
	return &SimpleTrunk{
		NewNode(0, true),
		0,
		numberOfFullList,
	}
}

/*
CachedTrunk is a Trunk with a LRU cache
*/
type CachedTrunk struct {
	*SimpleTrunk
	cache *lru.Cache
}

/*
NewCachedTrunk returns a new CachedTrunk
*/
func NewCachedTrunk(size int, numberOfFullList int) (*CachedTrunk, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &CachedTrunk{
		NewTrunk(numberOfFullList),
		cache,
	}, nil
}

func (t *SimpleTrunk) Append(nm *net.IPNet, data interface{}) {
	ones, _ := nm.Mask.Size()
	node := t.Node
	for i := 0; i < ones/8; i++ {
		node = node.SonOrNew(nm.IP[i], i < t.numberOfFullList)
	}
	node.Leafs = append(node.Leafs, &Leaf{
		Netmask: nm,
		Data:    data,
	})
	t.size++
}

func (t *SimpleTrunk) Size() int {
	return t.size
}

func (t *SimpleTrunk) Get(ip net.IP) (interface{}, bool) {
	ip = ip.To4()
	node := t.Node
	var n *Node
	for i := 0; i < 4; i++ {
		n = node.Son(ip[i])
		if n == nil {
			return nil, false
		}
		for _, leaf := range n.Leafs {
			if leaf.Netmask.Contains(ip) {
				return leaf.Data, true
			}
		}
		node = n
	}
	return nil, false
}

func (c *CachedTrunk) Get(ip net.IP) (interface{}, bool) {
	key := binary.BigEndian.Uint32(ip.To4())
	v, ok := c.cache.Get(key)
	if ok {
		vv := v.(response)
		return vv.value, vv.ok
	}
	value, ok := c.SimpleTrunk.Get(ip)
	c.cache.Add(key, response{ok, value})
	return value, ok
}

// Leaf is the end of a tree path
type Leaf struct {
	Netmask *net.IPNet
	Data    interface{}
}

type response struct {
	ok    bool
	value interface{}
}

func (t *SimpleTrunk) Dump(w io.Writer) {
	dump(w, 0, t.Node)
}

func dump(w io.Writer, tabs int, node *Node) {
	for _, son := range node.Sons {
		for i := 0; i < tabs; i++ {
			fmt.Fprint(w, "-")
		}
		fmt.Fprintf(w, "%x", son.Name)
		for _, leaf := range son.Leafs {
			fmt.Fprintf(w, " %v", leaf.Netmask)
		}
		fmt.Fprint(w, "\n")
		dump(w, tabs+1, son)
	}
}
