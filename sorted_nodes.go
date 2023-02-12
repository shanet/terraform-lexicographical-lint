package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jinzhu/copier"
)

const SORT_NAME = 1
const SORT_LINE_NUMBER = 2

const NODE_ATTRIBUTE = 1
const NODE_BLOCK = 2

type SortedNodes struct {
	Nodes  []SortedNode
	sortBy int
}

func (this SortedNodes) SortBy(sortBy int) (SortedNodes SortedNodes) {
	// Make a deep copy before sorting
	copier.Copy(&SortedNodes.Nodes, &this.Nodes)

	SortedNodes.sortBy = sortBy
	sort.Sort(SortedNodes)

	return SortedNodes
}

// Implement the sort interface https://golang.org/pkg/sort/#Interface
// --------------------------------------------------------------------

func (this SortedNodes) Len() int {
	return len(this.Nodes)
}

func (this SortedNodes) Less(left int, right int) bool {
	if this.sortBy == SORT_NAME {
		return strings.ToLower(this.Nodes[left].Name) < strings.ToLower(this.Nodes[right].Name)
	}

	if this.sortBy == SORT_LINE_NUMBER {
		return this.Nodes[left].Line < this.Nodes[right].Line
	}

	return false
}

func (this SortedNodes) Swap(left int, right int) {
	tmpNode := this.Nodes[left]
	this.Nodes[left] = this.Nodes[right]
	this.Nodes[right] = tmpNode
}

// Various utility functions
// --------------------------------------------------------------------

func (this *SortedNodes) Append(node SortedNode) []SortedNode {
	this.Nodes = append(this.Nodes, node)
	return this.Nodes
}

func (this SortedNodes) DebugPrint(params ...string) {
	if len(params) > 0 {
		fmt.Printf("%s\n", params[0])
	}

	for _, node := range this.Nodes {
		fmt.Printf("%d: %s (%d)\n", node.Line, node.Name, node.Type)
	}
	fmt.Println()
}
