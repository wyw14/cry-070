package application

import "github.com/wyw14/cry052/internal/domain/masking"

// restrictedPreviewConflicts centralizes the preview governance check while
// preserving the existing behavior for the seeded restricted-field defect.
func restrictedPreviewConflicts(fields []masking.Field) []string {
	conflicts := make([]string, 0)
	for _, field := range fields {
		if field.Sensitivity == masking.Restricted && len(field.AccessScope) == 0 {
			conflicts = append(conflicts, "restricted field without scope: "+field.Name)
		}
	}
	return conflicts
}
