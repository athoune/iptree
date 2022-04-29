package tree

import "sort"

// Node describes a tree graph
type Node struct {
	Name     byte
	Sons     Nodes
	Leafs    []*Leaf
	fullList bool
}

// Nodes is a list of Node, Nodes is sortable
type Nodes []*Node

// Len is the nimber of Node
func (n Nodes) Len() int { return len(n) }

// Swap two Nodes
func (n Nodes) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

// Less
func (n Nodes) Less(i, j int) bool { return n[i].Name < n[j].Name }

// Son of a Node of this key
func (n *Node) Son(a byte) *Node {
	if n.fullList {
		return n.Sons[a]
	}
	if len(n.Sons) == 0 {
		return nil
	}
	r := sort.Search(len(n.Sons), func(i int) bool {
		return n.Sons[i].Name >= a
	})
	if r < len(n.Sons) {
		return n.Sons[r]
	}
	return nil
}

// SonOrNew return the next Node in the graph
func (n *Node) SonOrNew(a byte, full bool) *Node {
	node := n.Son(a)
	if node != nil {
		return node
	}
	node = NewNode(a, full)
	if n.fullList {
		n.Sons[a] = node
	} else {
		n.Sons = append(n.Sons, node)
		sort.Sort(n.Sons)
	}
	return node
}

// NewNode return a new Node
func NewNode(name byte, full bool) *Node {
	var size int
	if full {
		size = 256
	} else {
		size = 0
	}
	return &Node{
		Name:     name,
		Sons:     make(Nodes, size),
		Leafs:    make([]*Leaf, 0),
		fullList: full,
	}
}
