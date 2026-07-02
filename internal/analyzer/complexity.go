package analyzer

import (
	sitter "github.com/smacker/go-tree-sitter"
)

// computeCyclomaticComplexity computes McCabe's cyclomatic complexity
// CC = number of decision points + 1
func computeCyclomaticComplexity(node *sitter.Node, source []byte) int {
	cc := 1 // base complexity

	var walk func(n *sitter.Node, isTopLevel bool)
	walk = func(n *sitter.Node, isTopLevel bool) {
		if n == nil {
			return
		}

		typeName := n.Type()

		// Don't recurse into nested functions (they have their own CC)
		// But allow traversal of the top-level function node
		if !isTopLevel && isNestedFunction(n) {
			return
		}

		// Decision points that increase CC
		switch typeName {
		case "if_statement", "elif_clause", "else_clause",
			"for_statement", "while_statement", "do_statement",
			"try_statement", "except_clause", "catch_clause",
			"switch_statement", "case_statement", "match_statement", "match_arm",
			"if_expression", "match_expression":
			cc++
		}

		// Logical operators also count as decision points
		switch typeName {
		case "and", "or", "&&", "||":
			cc++
		case "binary_operator":
			content := n.Content(source)
			if content == "&&" || content == "||" || content == "and" || content == "or" {
				cc++
			}
		case "comparison_operator", "conditional_expression":
			cc++
		}

		for i := 0; i < int(n.ChildCount()); i++ {
			walk(n.Child(i), false)
		}
	}

	walk(node, true)
	return cc
}

// computeCognitiveComplexity computes cognitive complexity
// Based on SonarSource's cognitive complexity specification
func computeCognitiveComplexity(node *sitter.Node, nestingLevel int, source []byte) int {
	total := 0

	var walk func(n *sitter.Node, level int, isTopLevel bool)
	walk = func(n *sitter.Node, level int, isTopLevel bool) {
		if n == nil {
			return
		}

		typeName := n.Type()

		// Don't recurse into nested functions
		if !isTopLevel && isNestedFunction(n) {
			return
		}

		increment := 0

		// Control flow structures add increment + nesting bonus
		switch typeName {
		case "if_statement", "if_expression":
			increment = 1 + level
		case "elif_clause":
			increment = 1
		case "else_clause":
			increment = 1
		case "for_statement", "while_statement", "do_statement":
			increment = 1 + level
		case "try_statement":
			increment = 1
		case "except_clause", "catch_clause":
			increment = 1
		case "switch_statement", "match_statement":
			increment = 1 + level
		case "case_statement", "match_arm":
			increment = 1
		}

		// Logical operators
		switch typeName {
		case "and", "or", "&&", "||":
			increment++
		case "binary_operator":
			content := n.Content(source)
			if content == "&&" || content == "||" || content == "and" || content == "or" {
				increment++
			}
		}

		total += increment

		// Nesting structures increase nesting level for children
		if isNestingStructure(n) {
			newLevel := level + 1
			for i := 0; i < int(n.ChildCount()); i++ {
				walk(n.Child(i), newLevel, false)
			}
			return
		}

		for i := 0; i < int(n.ChildCount()); i++ {
			walk(n.Child(i), level, false)
		}
	}

	walk(node, nestingLevel, true)
	return total
}

// isNestedFunction checks if a node is a nested function definition
func isNestedFunction(node *sitter.Node) bool {
	typeName := node.Type()
	return typeName == "function_definition" ||
		typeName == "function_declaration" ||
		typeName == "function" ||
		typeName == "arrow_function" ||
		typeName == "method_definition" ||
		typeName == "generator_function" ||
		typeName == "function_item" ||
		typeName == "lambda"
}

// isNestingStructure checks if a node increases nesting level
func isNestingStructure(node *sitter.Node) bool {
	typeName := node.Type()
	return typeName == "if_statement" ||
		typeName == "for_statement" ||
		typeName == "while_statement" ||
		typeName == "do_statement" ||
		typeName == "switch_statement" ||
		typeName == "match_statement" ||
		typeName == "try_statement" ||
		typeName == "if_expression" ||
		typeName == "match_expression"
}
