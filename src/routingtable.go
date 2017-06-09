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
				fmt.Println(prefix.String(), prefix.Size, deep)
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
func (t *Table) Find(n Node) {}

// Save Table
func (t *Table) Save() {}

// Load Table
func (t *Table) Load(n Node) {}

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
