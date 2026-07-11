package analyzer

import (
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/EdgarOrtegaRamirez/codemetrics/internal/models"
	"github.com/EdgarOrtegaRamirez/codemetrics/internal/parser"
)

// FindHotspots analyzes a project and returns the top N hotspot functions.
// A hotspot is a function with high composite complexity score, combining
// cyclomatic complexity, cognitive complexity, nesting depth, and function length.
func (a *Analyzer) FindHotspots(rootPath string, topN int) (*models.HotspotReport, error) {
	report := &models.HotspotReport{
		RootPath:    rootPath,
		GeneratedAt: time.Now(),
		TopCount:    topN,
	}

	var allFuncs []models.HotspotResult
	totalFiles := 0

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

		totalFiles++

		for _, fn := range fm.Functions {
			score := computeCompositeScore(fn)
			relPath := strings.TrimPrefix(path, rootPath)
			relPath = strings.TrimPrefix(relPath, "/")

			allFuncs = append(allFuncs, models.HotspotResult{
				FilePath:       relPath,
				FuncName:       fn.Name,
				Line:           fn.Line,
				Cyclomatic:     fn.Cyclomatic,
				Cognitive:      fn.Cognitive,
				NestingDepth:   fn.NestingDepth,
				LinesOfCode:    fn.LinesOfCode,
				ParameterCount: fn.ParameterCount,
				CompositeScore: score,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by composite score descending
	sort.Slice(allFuncs, func(i, j int) bool {
		return allFuncs[i].CompositeScore > allFuncs[j].CompositeScore
	})

	report.TotalFiles = totalFiles
	report.TotalFuncs = len(allFuncs)

	// Take top N
	if topN <= 0 || topN > len(allFuncs) {
		topN = len(allFuncs)
	}
	report.Hotspots = allFuncs[:topN]

	return report, nil
}

// computeCompositeScore computes a weighted composite complexity score.
// Uses normalized values for each metric (0-1 range) and combines them.
// Higher score = more likely to be a bug hotspot.
//
// Weights:
//   - Cyclomatic complexity: 0.35 (most predictive of defects per research)
//   - Cognitive complexity: 0.25 (readability/maintainability)
//   - Nesting depth: 0.20 (deep nesting is hard to reason about)
//   - Function length: 0.20 (long functions tend to do too much)
func computeCompositeScore(fn models.FuncMetrics) float64 {
	// Normalize using reasonable baselines
	cyclomaticNorm := normalize(float64(fn.Cyclomatic), 1.0, 50.0)
	cognitiveNorm := normalize(float64(fn.Cognitive), 1.0, 100.0)
	nestingNorm := normalize(float64(fn.NestingDepth), 1.0, 10.0)
	lengthNorm := normalize(float64(fn.LinesOfCode), 1.0, 200.0)

	score := cyclomaticNorm*0.35 + cognitiveNorm*0.25 + nestingNorm*0.20 + lengthNorm*0.20

	// Round to 2 decimal places
	return math.Round(score*100) / 100
}

// normalize normalizes a value to [0, 1] range using the formula:
// normalized = (value - min) / (max - min)
// Values above max are capped at 1.0.
func normalize(value, min, max float64) float64 {
	if value <= min {
		return 0.0
	}
	if value >= max {
		return 1.0
	}
	return (value - min) / (max - min)
}