package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/EdgarOrtegaRamirez/codemetrics/internal/models"
	"github.com/EdgarOrtegaRamirez/codemetrics/internal/parser"
)

// Analyzer computes code metrics using tree-sitter AST
type Analyzer struct {
	parser *parser.Parser
}

// New creates a new Analyzer
func New() *Analyzer {
	return &Analyzer{parser: parser.New()}
}

// AnalyzeFile computes metrics for a single file
func (a *Analyzer) AnalyzeFile(path string) (*models.FileMetrics, error) {
	root, lang, err := a.parser.ParseFile(path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	metrics := &models.FileMetrics{
		FilePath:   path,
		Language:   lang,
		AnalyzedAt: time.Now(),
	}

	// Compute LOC
	metrics.LinesOfCode = computeLOC(lines)

	// Compute function-level metrics
	functions := extractFunctions(root, lang, content)
	metrics.Functions = make([]models.FuncMetrics, 0, len(functions))

	for _, fn := range functions {
		fm := computeFuncMetrics(fn, lines, lang, content)
		metrics.Functions = append(metrics.Functions, fm)
	}

	// Compute aggregate complexity
	maxNesting := 0
	totalCC := 0
	totalCog := 0
	for _, fm := range metrics.Functions {
		totalCC += fm.Cyclomatic
		totalCog += fm.Cognitive
		if fm.NestingDepth > maxNesting {
			maxNesting = fm.NestingDepth
		}
	}

	metrics.TotalComplexity = models.Complexity{
		Cyclomatic: totalCC,
		Cognitive:  totalCog,
	}
	metrics.MaxNestingDepth = maxNesting

	return metrics, nil
}

// AnalyzeProject computes metrics for all supported files in a directory
func (a *Analyzer) AnalyzeProject(rootPath string) (*models.ProjectMetrics, error) {
	project := &models.ProjectMetrics{
		RootPath:   rootPath,
		AnalyzedAt: time.Now(),
	}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip files we can't read
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "__pycache__" || name == ".venv" || name == "target" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		lang := parser.DetectLanguage(path)
		if lang == models.LanguageUnknown {
			return nil
		}

		fm, err := a.AnalyzeFile(path)
		if err != nil {
			return nil // skip files that fail to parse
		}

		project.Files = append(project.Files, *fm)
		project.TotalFiles++
		project.TotalLOC.Total += fm.LinesOfCode.Total
		project.TotalLOC.Code += fm.LinesOfCode.Code
		project.TotalLOC.Comments += fm.LinesOfCode.Comments
		project.TotalLOC.Blanks += fm.LinesOfCode.Blanks
		project.TotalComplexity.Cyclomatic += fm.TotalComplexity.Cyclomatic
		project.TotalComplexity.Cognitive += fm.TotalComplexity.Cognitive

		if fm.TotalComplexity.Cyclomatic > project.MaxComplexity.Cyclomatic {
			project.MaxComplexity.Cyclomatic = fm.TotalComplexity.Cyclomatic
		}
		if fm.TotalComplexity.Cognitive > project.MaxComplexity.Cognitive {
			project.MaxComplexity.Cognitive = fm.TotalComplexity.Cognitive
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Compute summary
	project.Summary = computeSummary(project)

	return project, nil
}

// computeSummary computes summary statistics
func computeSummary(project *models.ProjectMetrics) models.Summary {
	s := models.Summary{}

	totalFuncs := 0
	totalFuncLen := 0
	maxFuncLen := 0
	totalCC := 0.0
	totalCog := 0.0
	complexFiles := 0

	for _, fm := range project.Files {
		if fm.TotalComplexity.Cyclomatic > 10 {
			complexFiles++
		}
		for _, fn := range fm.Functions {
			totalFuncs++
			totalFuncLen += fn.LinesOfCode
			if fn.LinesOfCode > maxFuncLen {
				maxFuncLen = fn.LinesOfCode
			}
			totalCC += float64(fn.Cyclomatic)
			totalCog += float64(fn.Cognitive)
		}
	}

	s.TotalFunctions = totalFuncs
	s.MaxFuncLength = maxFuncLen
	s.ComplexFiles = complexFiles

	if totalFuncs > 0 {
		s.AvgCyclomatic = totalCC / float64(totalFuncs)
		s.AvgCognitive = totalCog / float64(totalFuncs)
		s.AvgFuncLength = float64(totalFuncLen) / float64(totalFuncs)
	}

	return s
}
