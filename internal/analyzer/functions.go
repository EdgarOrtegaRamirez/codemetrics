package analyzer

import (
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/EdgarOrtegaRamirez/codemetrics/internal/models"
)

// funcNode represents an extracted function from the AST
type funcNode struct {
	name      string
	startLine int
	endLine   int
	node      *sitter.Node
}

// extractFunctions extracts all function definitions from the AST
func extractFunctions(root *sitter.Node, lang models.Language, source []byte) []funcNode {
	var funcs []funcNode

	var walk func(node *sitter.Node)
	walk = func(node *sitter.Node) {
		if node == nil {
			return
		}

		if isFunctionNode(node, lang) {
			fn := funcNode{
				name:      getFunctionName(node, lang, source),
				startLine: int(node.StartPoint().Row) + 1,
				endLine:   int(node.EndPoint().Row) + 1,
				node:      node,
			}
			funcs = append(funcs, fn)
		}

		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			walk(child)
		}
	}

	walk(root)
	return funcs
}

// isFunctionNode checks if a node is a function definition
func isFunctionNode(node *sitter.Node, lang models.Language) bool {
	typeName := node.Type()

	switch lang {
	case models.LanguagePython:
		return typeName == "function_definition" || typeName == "class_definition"
	case models.LanguageJavaScript, models.LanguageTypeScript:
		return typeName == "function_declaration" ||
			typeName == "function" ||
			typeName == "arrow_function" ||
			typeName == "method_definition" ||
			typeName == "generator_function"
	case models.LanguageGo:
		return typeName == "function_declaration" || typeName == "method_declaration"
	case models.LanguageRust:
		return typeName == "function_item"
	default:
		return false
	}
}

// getFunctionName extracts the name from a function node
func getFunctionName(node *sitter.Node, lang models.Language, source []byte) string {
	switch lang {
	case models.LanguagePython:
		// function_definition -> identifier
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "identifier" {
				return child.Content(source)
			}
		}
		return "<anonymous>"
	case models.LanguageJavaScript, models.LanguageTypeScript:
		// Look for name child
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "identifier" || child.Type() == "property_identifier" {
				return child.Content(source)
			}
		}
		return "<anonymous>"
	case models.LanguageGo:
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "identifier" {
				return child.Content(source)
			}
		}
		return "<anonymous>"
	case models.LanguageRust:
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "identifier" {
				return child.Content(source)
			}
		}
		return "<anonymous>"
	default:
		return "<unknown>"
	}
}

// computeFuncMetrics computes all metrics for a function
func computeFuncMetrics(fn funcNode, lines []string, lang models.Language, source []byte) models.FuncMetrics {
	fm := models.FuncMetrics{
		Name:    fn.name,
		Line:    fn.startLine,
		EndLine: fn.endLine,
	}

	// Lines of code for this function
	if fn.startLine-1 >= 0 && fn.endLine <= len(lines) {
		fm.LinesOfCode = fn.endLine - fn.startLine + 1
	}

	// Count parameters
	fm.ParameterCount = countParameters(fn.node, lang)

	// Compute nesting depth
	fm.NestingDepth = computeNestingDepth(fn.node)

	// Compute cyclomatic complexity
	fm.Cyclomatic = computeCyclomaticComplexity(fn.node, source)

	// Compute cognitive complexity
	fm.Cognitive = computeCognitiveComplexity(fn.node, 0, source)

	return fm
}

// countParameters counts the number of parameters in a function
func countParameters(node *sitter.Node, lang models.Language) int {
	count := 0

	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}

		switch lang {
		case models.LanguagePython:
			if n.Type() == "parameters" {
				for i := 0; i < int(n.ChildCount()); i++ {
					child := n.Child(i)
					if child.Type() == "identifier" || child.Type() == "default_parameter" || child.Type() == "typed_parameter" {
						count++
					}
				}
				return // don't recurse into parameters
			}
		case models.LanguageJavaScript, models.LanguageTypeScript:
			if n.Type() == "formal_parameters" {
				for i := 0; i < int(n.ChildCount()); i++ {
					child := n.Child(i)
					if child.Type() == "identifier" || child.Type() == "required_parameter" || child.Type() == "optional_parameter" || child.Type() == "rest_pattern" {
						count++
					}
				}
				return
			}
		case models.LanguageGo:
			if n.Type() == "parameter_list" {
				for i := 0; i < int(n.ChildCount()); i++ {
					child := n.Child(i)
					if child.Type() == "parameter_declaration" {
						// Count identifiers in parameter_declaration
						count += countIdentifiers(child)
					}
				}
				return
			}
		case models.LanguageRust:
			if n.Type() == "parameters" {
				for i := 0; i < int(n.ChildCount()); i++ {
					child := n.Child(i)
					if child.Type() == "parameter" {
						count++
					}
				}
				return
			}
		}

		for i := 0; i < int(n.ChildCount()); i++ {
			walk(n.Child(i))
		}
	}

	walk(node)
	return count
}

// countIdentifiers counts identifier nodes (for Go parameter counting)
func countIdentifiers(node *sitter.Node) int {
	count := 0
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "identifier" {
			count++
		}
	}
	if count == 0 {
		count = 1 // at least one parameter
	}
	return count
}

// computeNestingDepth computes the maximum nesting depth of control structures
func computeNestingDepth(node *sitter.Node) int {
	maxDepth := 0

	var walk func(n *sitter.Node, depth int)
	walk = func(n *sitter.Node, depth int) {
		if n == nil {
			return
		}

		if isControlStructure(n) {
			depth++
			if depth > maxDepth {
				maxDepth = depth
			}
		}

		for i := 0; i < int(n.ChildCount()); i++ {
			walk(n.Child(i), depth)
		}
	}

	walk(node, 0)
	return maxDepth
}

// isControlStructure checks if a node is a control flow structure
func isControlStructure(node *sitter.Node) bool {
	typeName := node.Type()
	switch typeName {
	case "if_statement", "elif_clause", "else_clause",
		"for_statement", "while_statement", "do_statement",
		"try_statement", "except_clause", "catch_clause",
		"switch_statement", "case_statement", "match_statement", "match_arm",
		"if_expression", "match_expression", "block":
		return true
	default:
		return false
	}
}
