# CodeMetrics

A fast, multi-language code complexity analyzer using tree-sitter AST parsing.

## Features

- **Multi-language support**: Python, JavaScript, TypeScript, Go, Rust
- **AST-based analysis**: Uses tree-sitter for accurate, language-aware parsing
- **Cyclomatic complexity**: McCabe's CC metric for control flow complexity
- **Cognitive complexity**: SonarSource's metric for human-readable complexity
- **Lines of code**: Total, code, comments, and blank line counts
- **Function metrics**: Parameter count, nesting depth, function length
- **Multiple output formats**: Text, JSON, Markdown
- **Violation detection**: Find functions exceeding complexity thresholds
- **CI integration**: Non-zero exit codes for threshold violations

## Quick Start

```bash
# Install
go install github.com/EdgarOrtegaRamirez/codemetrics/cmd/codemetrics@latest

# Analyze a project
codemetrics analyze ./src

# Verbose output with per-function details
codemetrics analyze -v ./src

# JSON output
codemetrics analyze -f json ./src > report.json

# Find complex functions
codemetrics violations --cc-threshold 15 ./src

# Markdown report
codemetrics analyze -f markdown ./src > report.md
```

## Complexity Metrics

### Cyclomatic Complexity (CC)
Measures the number of linearly independent paths through a program's source code.

| CC | Severity | Description |
|----|----------|-------------|
| 1-5 | Low | Simple, easy to test |
| 6-10 | Medium | Moderate complexity |
| 11-20 | High | Complex, hard to test |
| 21+ | Critical | Very complex, refactor needed |

### Cognitive Complexity
Measures how difficult code is for humans to understand. Unlike cyclomatic complexity, it penalizes nesting and breaking the flow.

## Supported Languages

| Language | Extensions | AST Grammar |
|----------|------------|-------------|
| Python | `.py` | tree-sitter-python |
| JavaScript | `.js`, `.jsx`, `.mjs`, `.cjs` | tree-sitter-javascript |
| TypeScript | `.ts`, `.tsx`, `.mts`, `.cts` | tree-sitter-typescript |
| Go | `.go` | tree-sitter-go |
| Rust | `.rs` | tree-sitter-rust |

## CLI Commands

### `analyze`
Analyze code complexity for files or directories.

```
codemetrics analyze [options] <path>

Options:
  --format, -f    Output format: text (default), json, markdown
  --verbose, -v   Show detailed per-function metrics
  --file, -o      Output to file instead of stdout
  --cc-threshold  Cyclomatic complexity threshold (default: 10)
```

### `violations`
Find functions exceeding complexity thresholds.

```
codemetrics violations [options] <path>

Options:
  --format, -f    Output format: text (default), json, markdown
  --cc-threshold  Cyclomatic complexity threshold (default: 10)
  --file, -o      Output to file instead of stdout
```

### `version`
Show version information.

## Architecture

```
codemetrics/
├── cmd/codemetrics/        # CLI entry point
│   └── main.go
├── internal/
│   ├── models/             # Data structures
│   │   └── models.go
│   ├── parser/             # Tree-sitter based parser
│   │   └── parser.go
│   ├── analyzer/           # Metric computation
│   │   ├── analyzer.go     # Main analysis engine
│   │   ├── complexity.go   # CC & Cognitive complexity
│   │   ├── functions.go    # Function extraction
│   │   └── loc.go          # Lines of code
│   └── reporter/           # Output formatting
│       └── reporter.go
├── tests/                  # Test suite
└── fixtures/               # Test fixtures
```

## Use Cases

- **CI/CD gates**: Fail builds when complexity exceeds thresholds
- **Code reviews**: Identify complex functions needing refactoring
- **Technical debt**: Track complexity trends over time
- **Documentation**: Generate complexity reports for stakeholders
- **AI coding agents**: Understand codebase complexity for better suggestions

## License

MIT
