package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EdgarOrtegaRamirez/codemetrics/internal/analyzer"
	"github.com/EdgarOrtegaRamirez/codemetrics/internal/models"
	"github.com/EdgarOrtegaRamirez/codemetrics/internal/parser"
	"github.com/EdgarOrtegaRamirez/codemetrics/internal/reporter"
)

// getFixturePath returns the absolute path to a fixture file
func getFixturePath(relative string) string {
	// Try from project root
	paths := []string{
		filepath.Join("..", "fixtures", relative),
		filepath.Join("fixtures", relative),
		filepath.Join("..", "..", "fixtures", relative),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// Fallback to relative
	return filepath.Join("..", "fixtures", relative)
}

// TestAnalyzePythonFile tests Python file analysis
func TestAnalyzePythonFile(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze Python file: %v", err)
	}

	if fm.Language != "python" {
		t.Errorf("Expected language 'python', got '%s'", fm.Language)
	}

	if fm.LinesOfCode.Total == 0 {
		t.Error("Expected non-zero total LOC")
	}

	if fm.LinesOfCode.Blanks == 0 {
		t.Error("Expected non-zero blank lines")
	}

	if len(fm.Functions) == 0 {
		t.Error("Expected at least one function")
	}

	// simple() should have CC=1
	for _, fn := range fm.Functions {
		if fn.Name == "simple" && fn.Cyclomatic != 1 {
			t.Errorf("Expected simple() CC=1, got %d", fn.Cyclomatic)
		}
	}

	// complex_func should have CC > 5
	for _, fn := range fm.Functions {
		if fn.Name == "complex_func" && fn.Cyclomatic <= 5 {
			t.Errorf("Expected complex_func() CC > 5, got %d", fn.Cyclomatic)
		}
	}

	t.Logf("Python file: %d LOC, %d functions, CC=%d, Cognitive=%d",
		fm.LinesOfCode.Total, len(fm.Functions),
		fm.TotalComplexity.Cyclomatic, fm.TotalComplexity.Cognitive)

	for _, fn := range fm.Functions {
		t.Logf("  %s: CC=%d Cog=%d Params=%d Lines=%d",
			fn.Name, fn.Cyclomatic, fn.Cognitive, fn.ParameterCount, fn.LinesOfCode)
	}
}

// TestAnalyzeGoFile tests Go file analysis
func TestAnalyzeGoFile(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("go/sample.go"))
	if err != nil {
		t.Fatalf("Failed to analyze Go file: %v", err)
	}

	if fm.Language != "go" {
		t.Errorf("Expected language 'go', got '%s'", fm.Language)
	}

	if len(fm.Functions) == 0 {
		t.Error("Expected at least one function")
	}

	for _, fn := range fm.Functions {
		t.Logf("  %s: CC=%d Cog=%d Params=%d Lines=%d",
			fn.Name, fn.Cyclomatic, fn.Cognitive, fn.ParameterCount, fn.LinesOfCode)
	}
}

// TestAnalyzeJavaScriptFile tests JavaScript file analysis
func TestAnalyzeJavaScriptFile(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("javascript/sample.js"))
	if err != nil {
		t.Fatalf("Failed to analyze JavaScript file: %v", err)
	}

	if fm.Language != "javascript" {
		t.Errorf("Expected language 'javascript', got '%s'", fm.Language)
	}

	if len(fm.Functions) == 0 {
		t.Error("Expected at least one function")
	}

	for _, fn := range fm.Functions {
		t.Logf("  %s: CC=%d Cog=%d Params=%d Lines=%d",
			fn.Name, fn.Cyclomatic, fn.Cognitive, fn.ParameterCount, fn.LinesOfCode)
	}
}

// TestAnalyzeTypeScriptFile tests TypeScript file analysis
func TestAnalyzeTypeScriptFile(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("typescript/sample.ts"))
	if err != nil {
		t.Fatalf("Failed to analyze TypeScript file: %v", err)
	}

	if fm.Language != "typescript" {
		t.Errorf("Expected language 'typescript', got '%s'", fm.Language)
	}

	if len(fm.Functions) == 0 {
		t.Error("Expected at least one function")
	}

	for _, fn := range fm.Functions {
		t.Logf("  %s: CC=%d Cog=%d Params=%d Lines=%d",
			fn.Name, fn.Cyclomatic, fn.Cognitive, fn.ParameterCount, fn.LinesOfCode)
	}
}

// TestAnalyzeRustFile tests Rust file analysis
func TestAnalyzeRustFile(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("rust/sample.rs"))
	if err != nil {
		t.Fatalf("Failed to analyze Rust file: %v", err)
	}

	if fm.Language != "rust" {
		t.Errorf("Expected language 'rust', got '%s'", fm.Language)
	}

	if len(fm.Functions) == 0 {
		t.Error("Expected at least one function")
	}

	for _, fn := range fm.Functions {
		t.Logf("  %s: CC=%d Cog=%d Params=%d Lines=%d",
			fn.Name, fn.Cyclomatic, fn.Cognitive, fn.ParameterCount, fn.LinesOfCode)
	}
}

// TestAnalyzeProject tests project-wide analysis
func TestAnalyzeProject(t *testing.T) {
	a := analyzer.New()
	pm, err := a.AnalyzeProject(getFixturePath(""))
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	if pm.TotalFiles == 0 {
		t.Error("Expected at least one file")
	}

	t.Logf("Project: %d files, %d LOC, %d functions, CC=%d, Cognitive=%d",
		pm.TotalFiles, pm.TotalLOC.Total, pm.Summary.TotalFunctions,
		pm.TotalComplexity.Cyclomatic, pm.TotalComplexity.Cognitive)
}

// TestComplexitySeverity tests complexity severity classification
func TestComplexitySeverity(t *testing.T) {
	tests := []struct {
		cc       int
		expected string
	}{
		{1, "low"},
		{5, "low"},
		{6, "medium"},
		{10, "medium"},
		{11, "high"},
		{20, "high"},
		{21, "critical"},
		{50, "critical"},
	}

	for _, tt := range tests {
		severity := string(models.GetSeverity(tt.cc))
		if severity != tt.expected {
			t.Errorf("CC=%d: expected severity '%s', got '%s'", tt.cc, tt.expected, severity)
		}
	}
}

// TestLanguageDetection tests language detection from file extensions
func TestLanguageDetection(t *testing.T) {
	tests := []struct {
		path     string
		expected models.Language
	}{
		{"test.py", models.LanguagePython},
		{"test.js", models.LanguageJavaScript},
		{"test.jsx", models.LanguageJavaScript},
		{"test.mjs", models.LanguageJavaScript},
		{"test.ts", models.LanguageTypeScript},
		{"test.tsx", models.LanguageTypeScript},
		{"test.go", models.LanguageGo},
		{"test.rs", models.LanguageRust},
		{"test.txt", models.LanguageUnknown},
		{"test.java", models.LanguageUnknown},
		{"test.rb", models.LanguageUnknown},
	}

	for _, tt := range tests {
		lang := parser.DetectLanguage(tt.path)
		if lang != tt.expected {
			t.Errorf("DetectLanguage(%s): expected %s, got %s", tt.path, tt.expected, lang)
		}
	}
}

// TestReporterFormats tests all reporter output formats
func TestReporterFormats(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	formats := []reporter.Format{
		reporter.FormatText,
		reporter.FormatJSON,
		reporter.FormatMarkdown,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			f, err := os.CreateTemp("", "report-*")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(f.Name())
			defer f.Close()

			r := reporter.New(f, format, false)
			if err := r.WriteFileMetrics(fm); err != nil {
				t.Errorf("Failed to write %s report: %v", format, err)
			}
		})
	}
}

// TestReporterVerbose tests verbose output
func TestReporterVerbose(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	f, err := os.CreateTemp("", "report-verbose-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	r := reporter.New(f, reporter.FormatText, true)
	if err := r.WriteFileMetrics(fm); err != nil {
		t.Errorf("Failed to write verbose report: %v", err)
	}

	// Check that file is not empty
	info, _ := f.Stat()
	if info.Size() == 0 {
		t.Error("Verbose report should not be empty")
	}
}

// TestReportViolations tests violation reporting
func TestReportViolations(t *testing.T) {
	violations := []models.ComplexityViolation{
		{
			FilePath:   "test.py",
			FuncName:   "complex_func",
			Line:       10,
			Cyclomatic: 15,
			Cognitive:  20,
			Severity:   "high",
		},
		{
			FilePath:   "test.go",
			FuncName:   "Process",
			Line:       50,
			Cyclomatic: 25,
			Cognitive:  30,
			Severity:   "critical",
		},
	}

	formats := []reporter.Format{
		reporter.FormatText,
		reporter.FormatJSON,
		reporter.FormatMarkdown,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			f, err := os.CreateTemp("", "violations-*")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(f.Name())
			defer f.Close()

			r := reporter.New(f, format, false)
			if err := r.WriteViolations(violations); err != nil {
				t.Errorf("Failed to write violations in %s format: %v", format, err)
			}
		})
	}
}

// TestReportNoViolations tests no violation case
func TestReportNoViolations(t *testing.T) {
	f, err := os.CreateTemp("", "no-violations-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	r := reporter.New(f, reporter.FormatText, false)
	if err := r.WriteViolations([]models.ComplexityViolation{}); err != nil {
		t.Errorf("Failed to write empty violations: %v", err)
	}

	// Check that file has content (the "no violations" message)
	info, _ := f.Stat()
	if info.Size() == 0 {
		t.Error("No-violations report should have a message")
	}
}

// TestReportProjectMarkdown tests project-level markdown report
func TestReportProjectMarkdown(t *testing.T) {
	a := analyzer.New()
	pm, err := a.AnalyzeProject(getFixturePath(""))
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	f, err := os.CreateTemp("", "project-md-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	r := reporter.New(f, reporter.FormatMarkdown, false)
	if err := r.WriteProjectMetrics(pm); err != nil {
		t.Errorf("Failed to write project markdown: %v", err)
	}

	info, _ := f.Stat()
	if info.Size() == 0 {
		t.Error("Project markdown report should not be empty")
	}
}

// TestFileMetricsStructure tests file metrics structure completeness
func TestFileMetricsStructure(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	// Verify all fields are populated
	if fm.FilePath == "" {
		t.Error("FilePath should not be empty")
	}

	if fm.Language == "" {
		t.Error("Language should not be empty")
	}

	if fm.LinesOfCode.Total == 0 {
		t.Error("Total LOC should not be zero")
	}

	if fm.LinesOfCode.Code == 0 {
		t.Error("Code LOC should not be zero")
	}

	if fm.AnalyzedAt.IsZero() {
		t.Error("AnalyzedAt should not be zero")
	}

	// Verify function metrics
	for _, fn := range fm.Functions {
		if fn.Name == "" {
			t.Error("Function name should not be empty")
		}

		if fn.Line == 0 {
			t.Errorf("Function %s: line should not be zero", fn.Name)
		}

		if fn.EndLine < fn.Line {
			t.Errorf("Function %s: endLine %d < line %d", fn.Name, fn.EndLine, fn.Line)
		}

		if fn.Cyclomatic < 1 {
			t.Errorf("Function %s: CC should be >= 1, got %d", fn.Name, fn.Cyclomatic)
		}

		if fn.Cognitive < 0 {
			t.Errorf("Function %s: Cognitive should be >= 0, got %d", fn.Name, fn.Cognitive)
		}
	}
}

// TestProjectMetricsStructure tests project metrics structure completeness
func TestProjectMetricsStructure(t *testing.T) {
	a := analyzer.New()
	pm, err := a.AnalyzeProject(getFixturePath(""))
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	if pm.RootPath == "" {
		t.Error("RootPath should not be empty")
	}

	if pm.TotalFiles == 0 {
		t.Error("TotalFiles should not be zero")
	}

	if pm.TotalLOC.Total == 0 {
		t.Error("TotalLOC.Total should not be zero")
	}

	if pm.TotalLOC.Code == 0 {
		t.Error("TotalLOC.Code should not be zero")
	}

	if pm.Summary.TotalFunctions == 0 {
		t.Error("Summary.TotalFunctions should not be zero")
	}

	if pm.AnalyzedAt.IsZero() {
		t.Error("AnalyzedAt should not be zero")
	}
}

// TestComplexCCFunctions tests that complex functions are detected
func TestComplexCCFunctions(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	foundComplex := false
	for _, fn := range fm.Functions {
		if fn.Cyclomatic > 5 {
			foundComplex = true
			t.Logf("Complex function found: %s (CC=%d)", fn.Name, fn.Cyclomatic)
		}
	}

	if !foundComplex {
		t.Error("Expected at least one complex function (CC > 5)")
	}
}

// TestMultipleLanguagesInProject tests project with multiple languages
func TestMultipleLanguagesInProject(t *testing.T) {
	a := analyzer.New()
	pm, err := a.AnalyzeProject(getFixturePath(""))
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	languages := make(map[models.Language]bool)
	for _, fm := range pm.Files {
		languages[fm.Language] = true
	}

	t.Logf("Languages found: %v", languages)

	if len(languages) < 2 {
		t.Errorf("Expected at least 2 languages, got %d", len(languages))
	}
}

// TestNestingDepth tests nesting depth calculation
func TestNestingDepth(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	foundNested := false
	for _, fn := range fm.Functions {
		if fn.NestingDepth > 0 {
			foundNested = true
			t.Logf("Nested function: %s (depth=%d)", fn.Name, fn.NestingDepth)
		}
	}

	if !foundNested {
		t.Error("Expected at least one function with nesting depth > 0")
	}
}

// TestParameterCount tests parameter counting
func TestParameterCount(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	// Calculator class methods should have parameters
	for _, fn := range fm.Functions {
		if fn.Name == "add" && fn.ParameterCount == 0 {
			t.Error("Expected add() to have parameters")
		}
		if fn.Name == "compute" && fn.ParameterCount == 0 {
			t.Error("Expected compute() to have parameters")
		}
	}
}

// TestEmptyFile tests handling of empty files
func TestEmptyFile(t *testing.T) {
	// Create a temporary empty Python file
	tmpFile, err := os.CreateTemp("", "empty-*.py")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	a := analyzer.New()
	fm, err := a.AnalyzeFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to analyze empty file: %v", err)
	}

	if fm.LinesOfCode.Total != 1 {
		t.Errorf("Expected 1 total LOC for empty file (newline), got %d", fm.LinesOfCode.Total)
	}

	if len(fm.Functions) != 0 {
		t.Errorf("Expected 0 functions for empty file, got %d", len(fm.Functions))
	}
}

// TestInvalidFile tests handling of non-existent files
func TestInvalidFile(t *testing.T) {
	a := analyzer.New()
	_, err := a.AnalyzeFile("/nonexistent/file.py")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestUnsupportedLanguage tests handling of unsupported file types
func TestUnsupportedLanguage(t *testing.T) {
	// Create a temporary file with unsupported extension
	tmpFile, err := os.CreateTemp("", "test-*.xyz")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("some content")
	tmpFile.Close()

	a := analyzer.New()
	_, err = a.AnalyzeFile(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for unsupported language")
	}
}

// TestSeverityLevels tests all severity levels
func TestSeverityLevels(t *testing.T) {
	tests := []struct {
		cc       int
		expected models.Severity
	}{
		{0, models.SeverityLow},
		{1, models.SeverityLow},
		{5, models.SeverityLow},
		{6, models.SeverityMedium},
		{10, models.SeverityMedium},
		{11, models.SeverityHigh},
		{20, models.SeverityHigh},
		{21, models.SeverityCritical},
		{100, models.SeverityCritical},
	}

	for _, tt := range tests {
		severity := models.GetSeverity(tt.cc)
		if severity != tt.expected {
			t.Errorf("CC=%d: expected %s, got %s", tt.cc, tt.expected, severity)
		}
	}
}

// TestLOCBreakdown tests LOC breakdown accuracy
func TestLOCBreakdown(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	// Verify LOC breakdown adds up
	totalFromParts := fm.LinesOfCode.Code + fm.LinesOfCode.Comments + fm.LinesOfCode.Blanks
	if totalFromParts != fm.LinesOfCode.Total {
		t.Errorf("LOC breakdown doesn't add up: %d + %d + %d = %d != %d",
			fm.LinesOfCode.Code, fm.LinesOfCode.Comments, fm.LinesOfCode.Blanks,
			totalFromParts, fm.LinesOfCode.Total)
	}
}

// TestCognitiveComplexityNonNegative verifies cognitive >= 0
func TestCognitiveComplexityNonNegative(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	for _, fn := range fm.Functions {
		if fn.Cognitive < 0 {
			t.Errorf("Function %s: cognitive complexity should be >= 0, got %d",
				fn.Name, fn.Cognitive)
		}
	}
}

// TestSupportedLanguages tests parser supported languages list
func TestSupportedLanguages(t *testing.T) {
	p := parser.New()
	langs := p.SupportedLanguages()

	if len(langs) < 5 {
		t.Errorf("Expected at least 5 supported languages, got %d", len(langs))
	}

	// Verify specific languages are supported
	langSet := make(map[models.Language]bool)
	for _, lang := range langs {
		langSet[lang] = true
	}

	required := []models.Language{
		models.LanguagePython,
		models.LanguageJavaScript,
		models.LanguageTypeScript,
		models.LanguageGo,
		models.LanguageRust,
	}

	for _, lang := range required {
		if !langSet[lang] {
			t.Errorf("Expected language %s to be supported", lang)
		}
	}
}

// TestIsSupported tests parser IsSupported function
func TestIsSupported(t *testing.T) {
	p := parser.New()

	supported := []string{
		"test.py",
		"test.js",
		"test.ts",
		"test.go",
		"test.rs",
	}

	unsupported := []string{
		"test.txt",
		"test.java",
		"test.rb",
		"test.php",
	}

	for _, path := range supported {
		if !p.IsSupported(path) {
			t.Errorf("Expected %s to be supported", path)
		}
	}

	for _, path := range unsupported {
		if p.IsSupported(path) {
			t.Errorf("Expected %s to be unsupported", path)
		}
	}
}

// TestReportOutputToFile tests writing reports to files
func TestReportOutputToFile(t *testing.T) {
	a := analyzer.New()
	fm, err := a.AnalyzeFile(getFixturePath("python/sample.py"))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	outputPath := filepath.Join(t.TempDir(), "report.json")
	f, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	r := reporter.New(f, reporter.FormatJSON, false)
	if err := r.WriteFileMetrics(fm); err != nil {
		t.Errorf("Failed to write JSON report: %v", err)
	}

	// Verify file was written
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Output file should not be empty")
	}
}

// TestFindHotspots tests hotspot detection
func TestFindHotspots(t *testing.T) {
	a := analyzer.New()
	hr, err := a.FindHotspots(getFixturePath(""), 5)
	if err != nil {
		t.Fatalf("Failed to find hotspots: %v", err)
	}

	if hr == nil {
		t.Fatal("Expected non-nil hotspot report")
	}

	if hr.TotalFiles == 0 {
		t.Error("Expected at least one file in project")
	}

	if hr.TotalFuncs == 0 {
		t.Error("Expected at least one function")
	}

	if len(hr.Hotspots) == 0 {
		t.Error("Expected at least one hotspot")
	}

	if len(hr.Hotspots) > 5 {
		t.Errorf("Expected at most 5 hotspots, got %d", len(hr.Hotspots))
	}

	if hr.RootPath == "" {
		t.Error("RootPath should not be empty")
	}

	if hr.GeneratedAt.IsZero() {
		t.Error("GeneratedAt should not be zero")
	}

	// Verify hotspots are sorted by score descending
	for i := 1; i < len(hr.Hotspots); i++ {
		if hr.Hotspots[i].CompositeScore > hr.Hotspots[i-1].CompositeScore {
			t.Errorf("Hotspots not sorted by score descending: %f > %f",
				hr.Hotspots[i].CompositeScore, hr.Hotspots[i-1].CompositeScore)
		}
	}

	t.Logf("Hotspots: %d files, %d total funcs, %d top results", hr.TotalFiles, hr.TotalFuncs, len(hr.Hotspots))
	for _, h := range hr.Hotspots {
		t.Logf("  %s/%s: score=%.2f CC=%d Cog=%d Nest=%d Lines=%d",
			h.FilePath, h.FuncName, h.CompositeScore, h.Cyclomatic, h.Cognitive, h.NestingDepth, h.LinesOfCode)
	}
}

// TestFindHotspotsEmptyDir tests hotspots on an empty directory
func TestFindHotspotsEmptyDir(t *testing.T) {
	emptyDir := t.TempDir()
	a := analyzer.New()
	hr, err := a.FindHotspots(emptyDir, 10)
	if err != nil {
		t.Fatalf("Failed to find hotspots in empty dir: %v", err)
	}

	if hr == nil {
		t.Fatal("Expected non-nil hotspot report")
	}

	if hr.TotalFiles != 0 {
		t.Errorf("Expected 0 files in empty dir, got %d", hr.TotalFiles)
	}

	if hr.TotalFuncs != 0 {
		t.Errorf("Expected 0 funcs in empty dir, got %d", hr.TotalFuncs)
	}

	if len(hr.Hotspots) != 0 {
		t.Errorf("Expected 0 hotspots in empty dir, got %d", len(hr.Hotspots))
	}
}

// TestHotspotCompositeScore tests the composite score computation
func TestHotspotCompositeScore(t *testing.T) {
	a := analyzer.New()
	hr, err := a.FindHotspots(getFixturePath("python"), 10)
	if err != nil {
		t.Fatalf("Failed to find hotspots: %v", err)
	}

	if hr == nil || len(hr.Hotspots) == 0 {
		t.Fatal("Expected at least one hotspot")
	}

	// Verify scores are in valid range
	for _, h := range hr.Hotspots {
		if h.CompositeScore < 0 {
			t.Errorf("CompositeScore should not be negative, got %f", h.CompositeScore)
		}
		if h.FilePath == "" {
			t.Error("Hotspot file path should not be empty")
		}
		if h.FuncName == "" {
			t.Error("Hotspot func name should not be empty")
		}
	}
}

// TestHotspotReportFormats tests hotspot reporter output formats
func TestHotspotReportFormats(t *testing.T) {
	a := analyzer.New()
	hr, err := a.FindHotspots(getFixturePath(""), 3)
	if err != nil {
		t.Fatalf("Failed to find hotspots: %v", err)
	}

	formats := []reporter.Format{
		reporter.FormatText,
		reporter.FormatJSON,
		reporter.FormatMarkdown,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			f, err := os.CreateTemp("", "hotspots-*")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(f.Name())
			defer f.Close()

			r := reporter.New(f, format, false)
			if err := r.WriteHotspots(hr); err != nil {
				t.Errorf("Failed to write hotspots in %s format: %v", format, err)
			}

			info, _ := f.Stat()
			if info.Size() == 0 {
				t.Errorf("Hotspot report in %s format should not be empty", format)
			}
		})
	}
}

// TestHotspotAllFunctions tests that all functions are ranked when topN > total
func TestHotspotAllFunctions(t *testing.T) {
	a := analyzer.New()
	hr, err := a.FindHotspots(getFixturePath("python"), 999)
	if err != nil {
		t.Fatalf("Failed to find hotspots: %v", err)
	}

	if hr.TotalFuncs != len(hr.Hotspots) {
		t.Errorf("Expected all functions when topN > total: %d total vs %d returned",
			hr.TotalFuncs, len(hr.Hotspots))
	}
}

// TestHotspotTopN tests that topN limits results correctly
func TestHotspotTopN(t *testing.T) {
	a := analyzer.New()

	// Get all hotspots
	allHr, _ := a.FindHotspots(getFixturePath("python"), 999)

	// Request top 3
	topHr, err := a.FindHotspots(getFixturePath("python"), 3)
	if err != nil {
		t.Fatalf("Failed to find top 3 hotspots: %v", err)
	}

	// If there are at least 3 functions, the top results should match
	if allHr.TotalFuncs >= 3 {
		if len(topHr.Hotspots) != 3 {
			t.Errorf("Expected 3 hotspots for --top 3, got %d", len(topHr.Hotspots))
		}

		// Check the top score is the same as the top of all
		if topHr.Hotspots[0].CompositeScore != allHr.Hotspots[0].CompositeScore {
			t.Errorf("Top hotspot score should match: got %f, expected %f",
				topHr.Hotspots[0].CompositeScore, allHr.Hotspots[0].CompositeScore)
		}
	}
}
