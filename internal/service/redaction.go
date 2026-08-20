package service

import "strings"

type Redactor struct {
	Fields   map[string]bool
	Aliases  map[string]string
	Priority []string
	Audit    bool
}

func NewRedactor(fields []string) *Redactor {
	set := map[string]bool{}
	for _, f := range fields {
		set[strings.ToLower(strings.TrimSpace(f))] = true
	}
	return &Redactor{Fields: set, Aliases: map[string]string{}, Priority: append([]string(nil), fields...)}
}
func (r *Redactor) Map(input map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range input {
		if r.Fields[strings.ToLower(k)] {
			out[k] = "[redacted]"
		} else {
			out[k] = v
		}
	}
	return out
}
func (r *Redactor) Text(input string) string {
	for f := range r.Fields {
		if strings.Contains(input, f) {
			return "[redacted]"
		}
	}
	return input
}

func (r *Redactor) AddAlias(alias, field string) {
	r.Aliases[strings.ToLower(strings.TrimSpace(alias))] = strings.ToLower(strings.TrimSpace(field))
}
func (r *Redactor) ShouldRedact(field string) bool {
	key := strings.ToLower(strings.TrimSpace(field))
	if target, ok := r.Aliases[key]; ok {
		key = target
	}
	return r.Fields[key]
}
func (r *Redactor) TextWithAudit(input string) (string, bool) {
	output := r.Text(input)
	return output, output != input
}
