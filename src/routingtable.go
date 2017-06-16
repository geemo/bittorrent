package dht

import (
	"fmt"
	"strconv"
	"time"
)

type bucket struct {
	nodes       []Node
	lastChanged time.Time
}

// TableNode of binary prefix tree
type TableNode struct {
	left   *TableNode
	right  *TableNode
	bucket *bucket
}

// Table struct
type Table struct {
	root     *TableNode
	selfNode Node
	k        int
}

var (
	leftValue  = 1
	rightValue = 0
	maxDeep    = 160
)

func (tn *TableNode) isEnd() bool {
	return tn.bucket != nil
}

func (tn *TableNode) addNode(n Node) {
	tn.bucket.nodes = append(tn.bucket.nodes, n)
}

func (tn *TableNode) clear() {
	tn.bucket = nil
}

// NewTable return a new Table
func NewTable(root *TableNode, selfNode Node) *Table {
	return &Table{root, selfNode, 8}
}

// Insert TableNode
func (t *Table) Insert(n Node) {
	t.insert(t.root, n, 0)
}

func (t *Table) insert(tn *TableNode, n Node, deep int) {
	if deep >= maxDeep {
		return
	}

	if tn.isEnd() {
		if len(tn.bucket.nodes) < t.k {
			// has bucket and lenght smaller than k insert bucket
			tn.addNode(n)
		} else {
			prefix := t.selfNode.id.Prefix(n.id)
			if prefix.Size >= deep {
				nodes := tn.bucket.nodes
				tn.clear()
				for _, node := range nodes {
					t.insert(tn, node, deep)
				}
			}
		}
		return
	}

	bit := n.id.Bit(deep)
	var nodes []Node
	bucket := &bucket{append(nodes, n), time.Now()}
	newTn := &TableNode{nil, nil, bucket}

	if bit == leftValue {
		if tn.left != nil {
			t.insert(tn.left, n, deep+1)
		} else {
			tn.left = newTn
		}
	} else {
		if tn.right != nil {
			t.insert(tn.right, n, deep+1)
		} else {
			tn.right = newTn
		}
	}
	return
}

// Delete TableNode
func (t *Table) Delete(n Node) {}

// Find TableNode
func (t *Table) Find(n Node) (Node, bool) {
	nodes := t.Closest(n)
	for _, node := range nodes {
		if node.id.RawString() == n.id.RawString() {
			return node, true
		}
	}
	return Node{}, false
}

// Save Table
func (t *Table) Save() {}

// Load Table
func (t *Table) Load(n Node) {}

// Closest 获得与 n 相比最近 k 个 node
func (t *Table) Closest(n Node) []Node {
	return t.closest(t.root, n, []path{}, []Node{}, 0)
}

type path struct {
	tn   *TableNode
	deep int
}

func popPath(slice []path) (path, []path) {
	lenght := len(slice)
	n := slice[lenght-1]
	newSlice := slice[:lenght-1]
	return n, newSlice
}

func (t *Table) closest(tn *TableNode, n Node, p []path, nodes []Node, deep int) []Node {
	// 有 bucket
	if tn.bucket != nil {
		nodes = append(nodes, tn.bucket.nodes...)
	}

	// 拿到 k 个 node 返回
	if len(nodes) >= t.k {
		return nodes[:t.k]
	}

	bit := n.id.Bit(deep)

	// 优先相同前缀
	if bit == leftValue && tn.left != nil {
		if tn.right != nil {
			p = append(p, path{tn.right, deep + 1})
		}
		return t.closest(tn.left, n, p, nodes, deep+1)
	}

	if bit == rightValue && tn.right != nil {
		if tn.left != nil {
			p = append(p, path{tn.left, deep + 1})
		}
		return t.closest(tn.right, n, p, nodes, deep+1)
	}

	// 没有相同前缀，随便挑选分支
	if tn.left != nil {
		if tn.right != nil {
			p = append(p, path{tn.right, deep + 1})
		}
		return t.closest(tn.left, n, p, nodes, deep+1)
	}

	if tn.right != nil {
		if tn.left != nil {
			p = append(p, path{tn.left, deep + 1})
		}
		return t.closest(tn.right, n, p, nodes, deep+1)
	}

	// 没有可选分支， 返回之前分支
	if len(p) > 0 {
		onePath, newPath := popPath(p)
		return t.closest(onePath.tn, n, newPath, nodes, onePath.deep)
	}

	// 什么都没有，就返回 nodes
	return nodes
}

// Print table
func (t *Table) Print() {
	t.print(t.root, "")
}

func (t *Table) print(tb *TableNode, prefix string) {
	if tb.bucket != nil {
		fmt.Println("** prefix:", prefix, "**")
		for _, node := range tb.bucket.nodes {
			fmt.Println(IDToHex(node.id), node.ip, node.port)
		}
	}
	if tb.left != nil {
		t.print(tb.left, prefix+strconv.Itoa(leftValue))
	}
	if tb.right != nil {
		t.print(tb.right, prefix+strconv.Itoa(rightValue))
	}
}
