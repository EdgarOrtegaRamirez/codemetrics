package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/EdgarOrtegaRamirez/codemetrics/internal/models"
)

// Format represents an output format
type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
)

// Report generates reports from project metrics
type Report struct {
	format  Format
	writer  io.Writer
	verbose bool
}

// New creates a new Report writer
func New(w io.Writer, format Format, verbose bool) *Report {
	return &Report{format: format, writer: w, verbose: verbose}
}

// WriteProjectMetrics writes project-level metrics
func (r *Report) WriteProjectMetrics(pm *models.ProjectMetrics) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(pm)
	case FormatMarkdown:
		return r.writeMarkdown(pm)
	default:
		return r.writeText(pm)
	}
}

// WriteFileMetrics writes file-level metrics
func (r *Report) WriteFileMetrics(fm *models.FileMetrics) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(fm)
	case FormatMarkdown:
		return r.writeFileMarkdown(fm)
	default:
		return r.writeFileText(fm)
	}
}

// WriteViolations writes complexity violations
func (r *Report) WriteViolations(violations []models.ComplexityViolation) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(violations)
	case FormatMarkdown:
		return r.writeViolationsMarkdown(violations)
	default:
		return r.writeViolationsText(violations)
	}
}

func (r *Report) writeJSON(data interface{}) error {
	enc := json.NewEncoder(r.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func (r *Report) writeText(pm *models.ProjectMetrics) error {
	w := tabwriter.NewWriter(r.writer, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "═══ CodeMetrics Report ═══\n\n")
	fmt.Fprintf(w, "Root: %s\n", pm.RootPath)
	fmt.Fprintf(w, "Files: %d\n", pm.TotalFiles)
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "── Lines of Code ──\n")
	fmt.Fprintf(w, "  Total:    %d\n", pm.TotalLOC.Total)
	fmt.Fprintf(w, "  Code:     %d\n", pm.TotalLOC.Code)
	fmt.Fprintf(w, "  Comments: %d\n", pm.TotalLOC.Comments)
	fmt.Fprintf(w, "  Blanks:   %d\n", pm.TotalLOC.Blanks)
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "── Complexity ──\n")
	fmt.Fprintf(w, "  Cyclomatic: %d (max: %d)\n", pm.TotalComplexity.Cyclomatic, pm.MaxComplexity.Cyclomatic)
	fmt.Fprintf(w, "  Cognitive:  %d (max: %d)\n", pm.TotalComplexity.Cognitive, pm.MaxComplexity.Cognitive)
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "── Summary ──\n")
	fmt.Fprintf(w, "  Functions:        %d\n", pm.Summary.TotalFunctions)
	fmt.Fprintf(w, "  Avg CC:           %.1f\n", pm.Summary.AvgCyclomatic)
	fmt.Fprintf(w, "  Avg Cognitive:    %.1f\n", pm.Summary.AvgCognitive)
	fmt.Fprintf(w, "  Avg Func Length:  %.1f lines\n", pm.Summary.AvgFuncLength)
	fmt.Fprintf(w, "  Max Func Length:  %d lines\n", pm.Summary.MaxFuncLength)
	fmt.Fprintf(w, "  Complex Files:    %d (CC > 10)\n", pm.Summary.ComplexFiles)

	if r.verbose && len(pm.Files) > 0 {
		fmt.Fprintf(w, "\n── File Details ──\n")
		for _, fm := range pm.Files {
			fmt.Fprintf(w, "\n  %s (%s)\n", fm.FilePath, fm.Language)
			fmt.Fprintf(w, "    LOC: %d (code: %d, comments: %d, blanks: %d)\n",
				fm.LinesOfCode.Total, fm.LinesOfCode.Code, fm.LinesOfCode.Comments, fm.LinesOfCode.Blanks)
			fmt.Fprintf(w, "    CC: %d  Cognitive: %d  Nesting: %d\n",
				fm.TotalComplexity.Cyclomatic, fm.TotalComplexity.Cognitive, fm.MaxNestingDepth)
			for _, fn := range fm.Functions {
				fmt.Fprintf(w, "      %s (line %d-%d): CC=%d Cog=%d Params=%d Lines=%d\n",
					fn.Name, fn.Line, fn.EndLine, fn.Cyclomatic, fn.Cognitive, fn.ParameterCount, fn.LinesOfCode)
			}
		}
	}

	return w.Flush()
}

func (r *Report) writeFileText(fm *models.FileMetrics) error {
	w := tabwriter.NewWriter(r.writer, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "── %s (%s) ──\n", fm.FilePath, fm.Language)
	fmt.Fprintf(w, "  LOC: %d (code: %d, comments: %d, blanks: %d)\n",
		fm.LinesOfCode.Total, fm.LinesOfCode.Code, fm.LinesOfCode.Comments, fm.LinesOfCode.Blanks)
	fmt.Fprintf(w, "  CC: %d  Cognitive: %d  Nesting: %d\n",
		fm.TotalComplexity.Cyclomatic, fm.TotalComplexity.Cognitive, fm.MaxNestingDepth)
	fmt.Fprintf(w, "  Functions: %d\n", len(fm.Functions))

	for _, fn := range fm.Functions {
		severity := models.GetSeverity(fn.Cyclomatic)
		fmt.Fprintf(w, "    %-30s line %4d-%-4d CC=%-3d Cog=%-3d Params=%-2d Lines=%-4d [%s]\n",
			fn.Name, fn.Line, fn.EndLine, fn.Cyclomatic, fn.Cognitive, fn.ParameterCount, fn.LinesOfCode, severity)
	}

	return w.Flush()
}

func (r *Report) writeViolationsText(violations []models.ComplexityViolation) error {
	if len(violations) == 0 {
		fmt.Fprintln(r.writer, "No complexity violations found.")
		return nil
	}

	w := tabwriter.NewWriter(r.writer, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "── Complexity Violations ──\n\n")
	fmt.Fprintf(w, "  %-40s %-30s %s %s %s\n", "File", "Function", "CC", "Cog", "Severity")
	fmt.Fprintf(w, "  %-40s %-30s %s %s %s\n", strings.Repeat("─", 40), strings.Repeat("─", 30), strings.Repeat("─", 4), strings.Repeat("─", 4), strings.Repeat("─", 10))

	for _, v := range violations {
		fmt.Fprintf(w, "  %-40s %-30s %-4d %-4d %s\n",
			v.FilePath, v.FuncName, v.Cyclomatic, v.Cognitive, v.Severity)
	}

	return w.Flush()
}

func (r *Report) writeMarkdown(pm *models.ProjectMetrics) error {
	fmt.Fprintf(r.writer, "# CodeMetrics Report\n\n")
	fmt.Fprintf(r.writer, "**Root:** `%s`  \n", pm.RootPath)
	fmt.Fprintf(r.writer, "**Files:** %d\n\n", pm.TotalFiles)

	fmt.Fprintf(r.writer, "## Lines of Code\n\n")
	fmt.Fprintf(r.writer, "| Metric | Count |\n")
	fmt.Fprintf(r.writer, "|--------|-------|\n")
	fmt.Fprintf(r.writer, "| Total | %d |\n", pm.TotalLOC.Total)
	fmt.Fprintf(r.writer, "| Code | %d |\n", pm.TotalLOC.Code)
	fmt.Fprintf(r.writer, "| Comments | %d |\n", pm.TotalLOC.Comments)
	fmt.Fprintf(r.writer, "| Blanks | %d |\n\n", pm.TotalLOC.Blanks)

	fmt.Fprintf(r.writer, "## Complexity\n\n")
	fmt.Fprintf(r.writer, "| Metric | Total | Max |\n")
	fmt.Fprintf(r.writer, "|--------|-------|-----|\n")
	fmt.Fprintf(r.writer, "| Cyclomatic | %d | %d |\n", pm.TotalComplexity.Cyclomatic, pm.MaxComplexity.Cyclomatic)
	fmt.Fprintf(r.writer, "| Cognitive | %d | %d |\n\n", pm.TotalComplexity.Cognitive, pm.MaxComplexity.Cognitive)

	fmt.Fprintf(r.writer, "## Summary\n\n")
	fmt.Fprintf(r.writer, "| Metric | Value |\n")
	fmt.Fprintf(r.writer, "|--------|-------|\n")
	fmt.Fprintf(r.writer, "| Functions | %d |\n", pm.Summary.TotalFunctions)
	fmt.Fprintf(r.writer, "| Avg Cyclomatic | %.1f |\n", pm.Summary.AvgCyclomatic)
	fmt.Fprintf(r.writer, "| Avg Cognitive | %.1f |\n", pm.Summary.AvgCognitive)
	fmt.Fprintf(r.writer, "| Avg Func Length | %.1f lines |\n", pm.Summary.AvgFuncLength)
	fmt.Fprintf(r.writer, "| Max Func Length | %d lines |\n", pm.Summary.MaxFuncLength)
	fmt.Fprintf(r.writer, "| Complex Files | %d |\n\n", pm.Summary.ComplexFiles)

	return nil
}

func (r *Report) writeFileMarkdown(fm *models.FileMetrics) error {
	fmt.Fprintf(r.writer, "## %s\n\n", fm.FilePath)
	fmt.Fprintf(r.writer, "**Language:** %s  \n", fm.Language)
	fmt.Fprintf(r.writer, "**LOC:** %d (code: %d, comments: %d, blanks: %d)  \n",
		fm.LinesOfCode.Total, fm.LinesOfCode.Code, fm.LinesOfCode.Comments, fm.LinesOfCode.Blanks)
	fmt.Fprintf(r.writer, "**CC:** %d  **Cognitive:** %d  **Nesting:** %d\n\n",
		fm.TotalComplexity.Cyclomatic, fm.TotalComplexity.Cognitive, fm.MaxNestingDepth)

	if len(fm.Functions) > 0 {
		fmt.Fprintf(r.writer, "### Functions\n\n")
		fmt.Fprintf(r.writer, "| Name | Lines | CC | Cognitive | Params | Length |\n")
		fmt.Fprintf(r.writer, "|------|-------|----|-----------|--------|--------|\n")
		for _, fn := range fm.Functions {
			fmt.Fprintf(r.writer, "| %s | %d-%d | %d | %d | %d | %d |\n",
				fn.Name, fn.Line, fn.EndLine, fn.Cyclomatic, fn.Cognitive, fn.ParameterCount, fn.LinesOfCode)
		}
	}

	return nil
}

func (r *Report) writeViolationsMarkdown(violations []models.ComplexityViolation) error {
	if len(violations) == 0 {
		fmt.Fprintln(r.writer, "No complexity violations found.")
		return nil
	}

	fmt.Fprintf(r.writer, "## Complexity Violations\n\n")
	fmt.Fprintf(r.writer, "| File | Function | CC | Cognitive | Severity |\n")
	fmt.Fprintf(r.writer, "|------|----------|----|-----------|----------|\n")
	for _, v := range violations {
		fmt.Fprintf(r.writer, "| %s | %s | %d | %d | %s |\n",
			v.FilePath, v.FuncName, v.Cyclomatic, v.Cognitive, v.Severity)
	}

	return nil
}
