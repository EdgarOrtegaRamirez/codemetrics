package analyzer

import (
	"strings"

	"github.com/EdgarOrtegaRamirez/codemetrics/internal/models"
)

// computeLOC computes lines-of-code breakdown
func computeLOC(lines []string) models.LOC {
	loc := models.LOC{Total: len(lines)}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			loc.Blanks++
		} else if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") {
			loc.Comments++
		} else {
			loc.Code++
		}
	}

	return loc
}
