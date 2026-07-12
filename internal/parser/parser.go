package parser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	golang "github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/rust"
	"github.com/smacker/go-tree-sitter/typescript/typescript"

	"github.com/EdgarOrtegaRamirez/codemetrics/internal/models"
)

// Parser wraps tree-sitter parsing for multiple languages
type Parser struct {
	parsers map[models.Language]*sitter.Parser
}

// New creates a new Parser with all supported language grammars
func New() *Parser {
	p := &Parser{
		parsers: make(map[models.Language]*sitter.Parser),
	}
	p.parsers[models.LanguagePython] = sitter.NewParser()
	p.parsers[models.LanguagePython].SetLanguage(python.GetLanguage())

	p.parsers[models.LanguageJavaScript] = sitter.NewParser()
	p.parsers[models.LanguageJavaScript].SetLanguage(javascript.GetLanguage())

	p.parsers[models.LanguageTypeScript] = sitter.NewParser()
	p.parsers[models.LanguageTypeScript].SetLanguage(typescript.GetLanguage())

	p.parsers[models.LanguageGo] = sitter.NewParser()
	p.parsers[models.LanguageGo].SetLanguage(golang.GetLanguage())

	p.parsers[models.LanguageRust] = sitter.NewParser()
	p.parsers[models.LanguageRust].SetLanguage(rust.GetLanguage())

	return p
}

// DetectLanguage determines the language from a file extension
func DetectLanguage(path string) models.Language {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".py":
		return models.LanguagePython
	case ".js", ".jsx", ".mjs", ".cjs":
		return models.LanguageJavaScript
	case ".ts", ".tsx", ".mts", ".cts":
		return models.LanguageTypeScript
	case ".go":
		return models.LanguageGo
	case ".rs":
		return models.LanguageRust
	default:
		return models.LanguageUnknown
	}
}

// ParseFile parses a file and returns the AST root node and language
func (p *Parser) ParseFile(path string) (*sitter.Node, models.Language, error) {
	lang := DetectLanguage(path)
	parser, ok := p.parsers[lang]
	if !ok {
		return nil, lang, fmt.Errorf("unsupported language for file: %s", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, lang, fmt.Errorf("reading file %s: %w", path, err)
	}

	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, lang, fmt.Errorf("parsing file %s: %w", path, err)
	}

	return tree.RootNode(), lang, nil
}

// SupportedLanguages returns the list of supported languages
func (p *Parser) SupportedLanguages() []models.Language {
	langs := make([]models.Language, 0, len(p.parsers))
	for lang := range p.parsers {
		langs = append(langs, lang)
	}
	return langs
}

// IsSupported checks if a file's language is supported
func (p *Parser) IsSupported(path string) bool {
	lang := DetectLanguage(path)
	_, ok := p.parsers[lang]
	return ok
}
