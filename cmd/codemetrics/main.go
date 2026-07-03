package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/codemetrics/internal/analyzer"
	"github.com/EdgarOrtegaRamirez/codemetrics/internal/models"
	"github.com/EdgarOrtegaRamirez/codemetrics/internal/reporter"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "analyze":
		cmdAnalyze(args)
	case "violations":
		cmdViolations(args)
	case "version":
		fmt.Printf("codemetrics %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`CodeMetrics - Multi-Language Code Complexity Analyzer

Usage:
  codemetrics <command> [options]

Commands:
  analyze     Analyze code complexity for files or directories
  violations  Find functions exceeding complexity thresholds
  version     Show version
  help        Show this help message

Analyze Options:
  --format, -f    Output format: text (default), json, markdown
  --verbose, -v   Show detailed per-function metrics
  --file, -o      Output to file instead of stdout
  --cc-threshold  Cyclomatic complexity threshold (default: 10)

Examples:
  codemetrics analyze ./src
  codemetrics analyze -f json -v ./src
  codemetrics violations --cc-threshold 15 ./src
  codemetrics analyze -f markdown ./src > report.md`)
}

func cmdAnalyze(args []string) {
	format := reporter.FormatText
	verbose := false
	threshold := 10
	outputFile := ""
	paths := []string{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--format", "-f":
			if i+1 < len(args) {
				i++
				format = reporter.Format(args[i])
			}
		case "--verbose", "-v":
			verbose = true
		case "--file", "-o":
			if i+1 < len(args) {
				i++
				outputFile = args[i]
			}
		case "--cc-threshold":
			if i+1 < len(args) {
				i++
				fmt.Sscanf(args[i], "%d", &threshold)
			}
		default:
			paths = append(paths, args[i])
		}
	}

	if len(paths) == 0 {
		paths = []string{"."}
	}

	// Open output
	var w *os.File
	var err error
	if outputFile != "" {
		w, err = os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer w.Close()
	} else {
		w = os.Stdout
	}

	a := analyzer.New()
	r := reporter.New(w, format, verbose)
	hasViolation := false

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if info.IsDir() {
			pm, err := a.AnalyzeProject(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", path, err)
				os.Exit(1)
			}
			if err := r.WriteProjectMetrics(pm); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing report: %v\n", err)
				os.Exit(1)
			}
			for _, fm := range pm.Files {
				for _, fn := range fm.Functions {
					if fn.Cyclomatic > threshold {
						hasViolation = true
					}
				}
			}
		} else {
			fm, err := a.AnalyzeFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", path, err)
				os.Exit(1)
			}
			if err := r.WriteFileMetrics(fm); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing report: %v\n", err)
				os.Exit(1)
			}
			for _, fn := range fm.Functions {
				if fn.Cyclomatic > threshold {
					hasViolation = true
				}
			}
		}
	}

	if hasViolation {
		os.Exit(1)
	}
}

func cmdViolations(args []string) {
	format := reporter.FormatText
	threshold := 10
	outputFile := ""
	paths := []string{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--format", "-f":
			if i+1 < len(args) {
				i++
				format = reporter.Format(args[i])
			}
		case "--file", "-o":
			if i+1 < len(args) {
				i++
				outputFile = args[i]
			}
		case "--cc-threshold":
			if i+1 < len(args) {
				i++
				fmt.Sscanf(args[i], "%d", &threshold)
			}
		default:
			paths = append(paths, args[i])
		}
	}

	if len(paths) == 0 {
		paths = []string{"."}
	}

	// Open output
	var w *os.File
	var err error
	if outputFile != "" {
		w, err = os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer w.Close()
	} else {
		w = os.Stdout
	}

	a := analyzer.New()
	r := reporter.New(w, format, false)

	var violations []models.ComplexityViolation

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if info.IsDir() {
			pm, err := a.AnalyzeProject(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", path, err)
				os.Exit(1)
			}
			for _, fm := range pm.Files {
				for _, fn := range fm.Functions {
					if fn.Cyclomatic > threshold {
						relPath := strings.TrimPrefix(fm.FilePath, path)
						relPath = strings.TrimPrefix(relPath, "/")
						violations = append(violations, models.ComplexityViolation{
							FilePath:   relPath,
							FuncName:   fn.Name,
							Line:       fn.Line,
							Cyclomatic: fn.Cyclomatic,
							Cognitive:  fn.Cognitive,
							Severity:   models.GetSeverity(fn.Cyclomatic),
						})
					}
				}
			}
		} else {
			fm, err := a.AnalyzeFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", path, err)
				os.Exit(1)
			}
			for _, fn := range fm.Functions {
				if fn.Cyclomatic > threshold {
					violations = append(violations, models.ComplexityViolation{
						FilePath:   fm.FilePath,
						FuncName:   fn.Name,
						Line:       fn.Line,
						Cyclomatic: fn.Cyclomatic,
						Cognitive:  fn.Cognitive,
						Severity:   models.GetSeverity(fn.Cyclomatic),
					})
				}
			}
		}
	}

	if err := r.WriteViolations(violations); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing report: %v\n", err)
		os.Exit(1)
	}

	if len(violations) > 0 {
		os.Exit(1) // non-zero exit if violations found
	}
}
