/*
  This program uses the Hasicorp HCL library to read Terraform files and perform additional validation beyond
  what the default Terraform formatter provides.

  Namely this involves ensuring the following:
    * Block attributes and nested block names are ordered
    * Variable declarations are ordered
    * `source` and `count` attributes are first when in a list of attributes and `provider` is last in a list of blocks

  The following resources are helpful:
    * https://github.com/hashicorp/hcl2/blob/master/hcl/hclsyntax/structure.go
    * https://github.com/hashicorp/hcl2/blob/master/hcl/hclsyntax/expression.go
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

var currentFile string
var success bool
var variableBlocks SortedBody

func main() {
	// Assume current working directory if no directory is given
	directory := "."

	if len(os.Args) == 2 {
		directory = os.Args[1]
	}

	success = true
	walkDirectory(directory)

	if success {
		fmt.Printf("%s: No lints found\n", os.Args[0])
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func walkDirectory(directory string) {
	entries, error := ioutil.ReadDir(directory)

	if error != nil {
		panic(error)
	}

	for _, entry := range entries {
		filename := entry.Name()
		path := filepath.Join(directory, filename)

		// Ignore .terraform folders so local tfstate files are not linted
		if filename == ".terraform" {
			continue
		} else if entry.IsDir() {
			walkDirectory(path)
		} else {
			processFile(path)
		}
	}
}

func processFile(path string) {
	// Update the global current file variable for printing it while walking the AST
	currentFile = path

	extension := filepath.Ext(path)

	if extension == ".tf" {
		fileContents, error := ioutil.ReadFile(path)

		if error != nil {
			panic(error)
		}

		// Parse and walk the AST
		file, diagnostics := hclsyntax.ParseConfig(fileContents, path, hcl.Pos{Line: 1, Column: 1})

		if diagnostics.HasErrors() {
			panic(diagnostics.Error())
		}

		// Initialize a new variable blocks struct for the file
		variableBlocks = SortedBody{Type: "variable", Line: 1, File: currentFile}

		hclsyntax.VisitAll(file.Body.(*hclsyntax.Body), walk)

		// After walking the whole file, verify that variable blocks at the top of the file are sorted
		if !verifyVariablesSorted(variableBlocks) {
			success = false
		}
	}
}

func walk(nextNode hclsyntax.Node) hcl.Diagnostics {
	switch node := nextNode.(type) {
	case hclsyntax.Blocks:
		for _, block := range node {
			// Collect all `variable` blocks to check their order
			if len(block.Labels) > 0 && block.Type == "variable" {
				variableBlocks.Append(SortedNode{Name: block.Labels[0], Line: block.TypeRange.Start.Line, Type: NODE_BLOCK})
			}

			if !verifyBlockSorted(block) {
				success = false
			}
		}

	case *hclsyntax.ObjectConsExpr:
		if !verifyExpressionSorted(node) {
			success = false
		}
	}

	return nil
}

func verifyBlockSorted(block *hclsyntax.Block) bool {
	// Construct some labels to display in error messages
	blockLine := block.TypeRange.Start.Line
	blockName := block.Type

	if len(block.Labels) > 0 {
		blockName += " " + block.Labels[len(block.Labels)-1]
	}

	sortedBody := SortedBody{Type: "Block", Name: blockName, Line: blockLine, File: currentFile}

	// Collect all of the attributes and nested blocks within the block into separate lists
	for _, attribute := range block.Body.Attributes {
		nodeType := NODE_ATTRIBUTE

		// If the attribute spans multiple lines, treat it as a nested block instead
		if attribute.SrcRange.Start.Line != attribute.SrcRange.End.Line {
			nodeType = NODE_BLOCK
		}

		sortedBody.Append(SortedNode{Name: attribute.Name, Line: attribute.SrcRange.Start.Line, Type: nodeType})
	}

	for _, block := range block.Body.Blocks {
		sortedBody.Append(SortedNode{Name: block.Type, Line: block.TypeRange.Start.Line, Type: NODE_BLOCK})
	}

	sortedBody.Sort()

	if !sortedBody.IsSorted() {
		sortedBody.PrintNonSorted()
		return false
	}

	return true
}

func verifyExpressionSorted(expression *hclsyntax.ObjectConsExpr) bool {
	attributes := SortedBody{Type: "Expression", Line: expression.SrcRange.Start.Line, File: currentFile}

	// Collect all of the keys within an expression (a map literal, for example)
	for _, item := range expression.Items {
		switch key := item.KeyExpr.(type) {
		case *hclsyntax.ObjectConsKeyExpr:
			keyword := hcl.ExprAsKeyword(key.Wrapped)

			if len(keyword) > 0 {
				nodeType := NODE_ATTRIBUTE

				// If the expression value spans multiple lines, treat it as a nested block instead
				switch value := item.ValueExpr.(type) {
				case *hclsyntax.TemplateExpr:
					if value.SrcRange.Start.Line != value.SrcRange.End.Line {
						nodeType = NODE_BLOCK
					}
				}

				attributes.Append(SortedNode{Name: keyword, Line: key.StartRange().Start.Line, Type: nodeType})
			}
		}
	}

	attributes.Sort()

	// Verify that the expression keys are sorted
	if !attributes.IsSorted() {
		attributes.PrintNonSorted()
		return false
	}

	return true
}

func verifyVariablesSorted(sortedBody SortedBody) bool {
	sortedBody.Sort()

	if !sortedBody.IsSorted() {
		sortedBody.PrintNonSorted()
		return false
	}

	return true
}
