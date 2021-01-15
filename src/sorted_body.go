package main

import (
	"fmt"
	"strings"
)

type SortedBody struct {
	Type  string
	Name  string
	Line  int
	File  string
	Nodes SortedNodes

	SpecialAttributes SortedNodes
	Attributes        SortedNodes
	Blocks            SortedNodes
	SpecialBlocks     SortedNodes

	ExpectedNodes SortedNodes
}

func (this *SortedBody) Sort() {
	// Sort the current nodes by line number
	nodesSortedByLineNumber := this.Nodes.SortBy(SORT_LINE_NUMBER)

	// Separate the nodes into the different types
	for _, node := range nodesSortedByLineNumber.Nodes {
		if node.Type == NODE_ATTRIBUTE && node.IsSpecial() {
			this.SpecialAttributes.Append(node)
		}

		if node.Type == NODE_ATTRIBUTE && !node.IsSpecial() {
			this.Attributes.Append(node)
		}

		if node.Type == NODE_BLOCK && !node.IsSpecial() {
			this.Blocks.Append(node)
		}

		if node.Type == NODE_BLOCK && node.IsSpecial() {
			this.SpecialBlocks.Append(node)
		}
	}

	// Insert each type of node into a single list in the expected order
	this.insert(this.SpecialAttributes)
	this.insert(this.Attributes)
	this.insert(this.Blocks)
	this.insert(this.SpecialBlocks)
}

func (this *SortedBody) insert(nodes SortedNodes) {
	nodesSortedByName := nodes.SortBy(SORT_NAME)

	for _, node := range nodesSortedByName.Nodes {
		this.ExpectedNodes.Append(node)
	}
}

func (this SortedBody) IsSorted() bool {
	// Check that the expected sort order is the same as the actual sorting (sort by line number)
	nodesSortedByLineNumber := this.Nodes.SortBy(SORT_LINE_NUMBER)

	for index, node := range nodesSortedByLineNumber.Nodes {
		if node.Name != this.ExpectedNodes.Nodes[index].Name {
			return false
		}
	}

	return true
}

func (this SortedBody) PrintNonSorted() {
	location := fmt.Sprintf("%s:%d", this.File, this.Line)
	var formattedName string

	if len(this.Name) > 0 {
		formattedName = color(34, fmt.Sprintf(" \"%s\"", strings.Trim(this.Name, "\"")))
	}

	fmt.Printf("%s %s%s expected order:\n", color(33, location), this.Type, formattedName)

	nodesSortedByLineNumber := this.Nodes.SortBy(SORT_LINE_NUMBER)

	for index, node := range this.ExpectedNodes.Nodes {
		// Display the out of order nodes in a different color
		if node.Name != nodesSortedByLineNumber.Nodes[index].Name {
			fmt.Printf("    %s\n", color(31, node.Name))
		} else {
			fmt.Printf("    %s\n", node.Name)
		}
	}
}

func (this *SortedBody) Append(node SortedNode) []SortedNode {
	this.Nodes.Nodes = append(this.Nodes.Nodes, node)
	return this.Nodes.Nodes
}

func color(color int, message string) string {
	return fmt.Sprintf("\033[1;%dm%s\033[0m", color, message)
}
