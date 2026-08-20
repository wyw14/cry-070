package service

import "sort"

type Page struct{ Offset, Limit int }

func Normalize(p Page) Page {
	if p.Offset < 0 {
		p.Offset = 0
	}
	if p.Limit < 1 || p.Limit > 100 {
		p.Limit = 20
	}
	return p
}
func Strings(values []string, p Page) []string {
	p = Normalize(p)
	out := append([]string(nil), values...)
	sort.Strings(out)
	if p.Offset >= len(out) {
		return []string{}
	}
	end := p.Offset + p.Limit
	if end > len(out) {
		end = len(out)
	}
	return out[p.Offset:end]
}
