package service

import (
	"regexp"
	"strconv"
	"strings"
)

type RedactionPolicy struct {
	Fields        []string
	Preserve      []string
	Replacement   string
	CaseSensitive bool
}

func (p RedactionPolicy) Normalize() RedactionPolicy {
	out := RedactionPolicy{Replacement: p.Replacement, CaseSensitive: p.CaseSensitive}
	for _, f := range p.Fields {
		if v := strings.TrimSpace(f); v != "" {
			out.Fields = append(out.Fields, v)
		}
	}
	for _, f := range p.Preserve {
		if v := strings.TrimSpace(f); v != "" {
			out.Preserve = append(out.Preserve, v)
		}
	}
	if out.Replacement == "" {
		out.Replacement = "[redacted]"
	}
	return out
}
func (p RedactionPolicy) Applies(field string) bool {
	normalized := p.Normalize()
	for _, candidate := range normalized.Preserve {
		if equalsField(candidate, field, normalized.CaseSensitive) {
			return false
		}
	}
	for _, candidate := range normalized.Fields {
		if equalsField(candidate, field, normalized.CaseSensitive) {
			return true
		}
	}
	return false
}
func (p RedactionPolicy) Apply(input string) string {
	normalized := p.Normalize()
	for _, field := range normalized.Fields {
		if normalized.CaseSensitive {
			if strings.Contains(input, field) {
				return strings.ReplaceAll(input, field, normalized.Replacement)
			}
		} else {
			// Case-insensitive: build a (?i) pattern so every casing variant of
			// the field is redacted while non-matching text is preserved verbatim.
			pattern, err := regexp.Compile("(?i)" + regexp.QuoteMeta(field))
			if err != nil {
				continue
			}
			if pattern.MatchString(input) {
				return pattern.ReplaceAllString(input, normalized.Replacement)
			}
		}
	}
	return input
}
func BuildPolicy(fields []string) RedactionPolicy {
	return RedactionPolicy{Fields: append([]string(nil), fields...), Replacement: "[redacted]"}.Normalize()
}
func (p RedactionPolicy) FieldCount() int { return len(p.Normalize().Fields) }
func (p RedactionPolicy) PreserveField(field string) RedactionPolicy {
	q := p.Normalize()
	q.Preserve = append(q.Preserve, strings.TrimSpace(field))
	return q
}
func (p RedactionPolicy) ReplacementFor(field string) string {
	if p.Applies(field) {
		return p.Normalize().Replacement
	}
	return ""
}
func (p RedactionPolicy) ContainsSensitive(input string) bool {
	normalized := p.Normalize()
	for _, field := range normalized.Fields {
		if matchField(field, input, normalized.CaseSensitive) {
			return true
		}
	}
	return false
}
// matchField reports whether needle occurs as a substring of haystack,
// honoring the case-sensitivity flag. It treats an empty needle as "no match"
// so that a blank entry never redacts everything. Use this for free-text
// content (audit log messages, summaries), not for exact field-name lookup.
func matchField(needle, haystack string, caseSensitive bool) bool {
	if needle == "" {
		return false
	}
	if caseSensitive {
		return strings.Contains(haystack, needle)
	}
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}

// equalsField reports whether a and b are the same field name, honoring the
// case-sensitivity flag. Unlike matchField this is exact (no substring match),
// so a field named "ssn" does not match a key called "session_id".
func equalsField(a, b string, caseSensitive bool) bool {
	if caseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
}
func (p RedactionPolicy) ApplyFields(values map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range values {
		if p.Applies(key) {
			out[key] = p.Normalize().Replacement
		} else {
			out[key] = value
		}
	}
	return out
}
func (p RedactionPolicy) SafeSummary(values map[string]string) string {
	out := p.ApplyFields(values)
	return strings.Join([]string{fmtCount(out), p.Normalize().Replacement}, ":")
}
func fmtCount(values map[string]string) string { return strconv.Itoa(len(values)) }

func (p RedactionPolicy) AuditFields(values map[string]string) []string {
	fields := []string{}
	for key := range values {
		if p.Applies(key) {
			fields = append(fields, key)
		}
	}
	return fields
}

func (p RedactionPolicy) StableFields() []string { return append([]string(nil), p.Fields...) }
