package models

import "time"

// Language represents a supported programming language
type Language string

const (
	LanguagePython     Language = "python"
	LanguageJavaScript Language = "javascript"
	LanguageTypeScript Language = "typescript"
	LanguageGo         Language = "go"
	LanguageRust       Language = "rust"
	LanguageUnknown    Language = "unknown"
)

// FileMetrics holds metrics for a single file
type FileMetrics struct {
	FilePath        string       `json:"file_path"`
	Language        Language     `json:"language"`
	LinesOfCode     LOC          `json:"loc"`
	Functions       []FuncMetrics `json:"functions"`
	TotalComplexity Complexity   `json:"total_complexity"`
	MaxNestingDepth int          `json:"max_nesting_depth"`
	AnalyzedAt      time.Time    `json:"analyzed_at"`
}

// LOC holds lines-of-code breakdown
type LOC struct {
	Total    int `json:"total"`
	Code     int `json:"code"`
	Comments int `json:"comments"`
	Blanks   int `json:"blanks"`
}

// FuncMetrics holds metrics for a single function
type FuncMetrics struct {
	Name             string     `json:"name"`
	Line             int        `json:"line"`
	EndLine          int        `json:"end_line"`
	LinesOfCode      int        `json:"lines_of_code"`
	ParameterCount   int        `json:"parameter_count"`
	NestingDepth     int        `json:"nesting_depth"`
	Cyclomatic       int        `json:"cyclomatic_complexity"`
	Cognitive        int        `json:"cognitive_complexity"`
}

// Complexity holds aggregate complexity metrics
type Complexity struct {
	Cyclomatic int `json:"cyclomatic_complexity"`
	Cognitive  int `json:"cognitive_complexity"`
}

// ProjectMetrics holds metrics for an entire project/directory
type ProjectMetrics struct {
	RootPath        string        `json:"root_path"`
	TotalFiles      int           `json:"total_files"`
	TotalLOC        LOC           `json:"total_loc"`
	TotalComplexity Complexity    `json:"total_complexity"`
	MaxComplexity   Complexity    `json:"max_complexity"`
	Files           []FileMetrics `json:"files"`
	Summary         Summary       `json:"summary"`
	AnalyzedAt      time.Time     `json:"analyzed_at"`
}

// Summary provides a high-level summary
type Summary struct {
	AvgCyclomatic float64 `json:"avg_cyclomatic_complexity"`
	AvgCognitive  float64 `json:"avg_cognitive_complexity"`
	AvgFuncLength float64 `json:"avg_function_length"`
	MaxFuncLength int     `json:"max_function_length"`
	TotalFunctions int    `json:"total_functions"`
	ComplexFiles  int     `json:"complex_files"` // files with CC > 10
}

// Severity represents the severity of a complexity finding
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// ComplexityViolation represents a function that exceeds a complexity threshold
type ComplexityViolation struct {
	FilePath   string   `json:"file_path"`
	FuncName   string   `json:"func_name"`
	Line       int      `json:"line"`
	Cyclomatic int      `json:"cyclomatic_complexity"`
	Cognitive  int      `json:"cognitive_complexity"`
	Severity   Severity `json:"severity"`
}

// GetSeverity returns the severity for a given cyclomatic complexity value
func GetSeverity(cc int) Severity {
	switch {
	case cc <= 5:
		return SeverityLow
	case cc <= 10:
		return SeverityMedium
	case cc <= 20:
		return SeverityHigh
	default:
		return SeverityCritical
	}
}
