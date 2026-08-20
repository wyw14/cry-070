package service

import "strings"

type Redactor struct{ Fields map[string]bool }

func NewRedactor(fields []string) *Redactor {
	set := map[string]bool{}
	for _, f := range fields {
		set[strings.ToLower(strings.TrimSpace(f))] = true
	}
	return &Redactor{Fields: set}
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
		if strings.Contains(strings.ToLower(input), f) {
			return "[redacted]"
		}
	}
	return input
}
