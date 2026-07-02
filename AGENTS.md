# AGENTS.md — Notes for AI Agents

## Project Overview
CodeMetrics is a multi-language code complexity analyzer using tree-sitter AST parsing. It computes cyclomatic complexity, cognitive complexity, lines of code, and function-level metrics.

## Tech Stack
- **Language**: Go 1.24+
- **AST Parser**: tree-sitter (via smacker/go-tree-sitter)
- **Test**: Go standard testing
- **CLI**: Standard library flag package

## Project Structure
```
cmd/codemetrics/main.go    — CLI entry point
internal/
  models/models.go         — Data structures
  parser/parser.go         — Tree-sitter based parser
  analyzer/analyzer.go     — Main analysis engine
  analyzer/complexity.go   — CC & Cognitive complexity
  analyzer/functions.go    — Function extraction
  analyzer/loc.go          — Lines of code
  reporter/reporter.go     — Output formatting
tests/                     — Test files
fixtures/                  — Test fixtures (Python, JS, TS, Go, Rust)
```

## Build & Test
```bash
go build ./...            # Build all packages
go test ./... -v          # Run all tests
go vet ./...              # Run vet checks
```

## Key Concepts
- Uses tree-sitter for AST parsing (supports Python, JS, TS, Go, Rust)
- Cyclomatic complexity: counts decision points (if/for/while/switch/case)
- Cognitive complexity: penalizes nesting and flow breaks
- Source content must be passed to Content() calls for correct text extraction
- isNestedFunction check prevents recursing into nested function definitions

## CI
- GitHub Actions in `.github/workflows/ci.yml`
- Runs on push/PR to main
- Go 1.24 matrix

## Commit Convention
- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation
- `test:` — Tests
- `refactor:` — Code refactoring
- `chore:` — Maintenance
