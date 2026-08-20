package masking

import "strings"

func ValidatePipeline(p Pipeline) []string {
	out := []string{}
	if strings.TrimSpace(p.ID) == "" {
		out = append(out, "id required")
	}
	if strings.TrimSpace(p.Name) == "" {
		out = append(out, "name required")
	}
	if p.SourceID == "" {
		out = append(out, "source required")
	}
	return out
}
func ValidateField(f Field) []string {
	out := []string{}
	if f.ID == "" {
		out = append(out, "id required")
	}
	if f.Table == "" {
		out = append(out, "table required")
	}
	if f.Name == "" {
		out = append(out, "name required")
	}
	return out
}
