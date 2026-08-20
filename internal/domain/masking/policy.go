package masking

import "strings"

type Policy struct {
	Name, Owner                     string
	Allowed                         []Sensitivity
	RequirePreview, RequireApproval bool
	Version                         int
}

func (p Policy) Allows(level Sensitivity) bool {
	for _, v := range p.Allowed {
		if v == level {
			return true
		}
	}
	return false
}
func (p Policy) Validate() []string {
	out := []string{}
	if strings.TrimSpace(p.Name) == "" {
		out = append(out, "name required")
	}
	if strings.TrimSpace(p.Owner) == "" {
		out = append(out, "owner required")
	}
	if len(p.Allowed) == 0 {
		out = append(out, "sensitivity scope required")
	}
	if p.Version < 1 {
		out = append(out, "version must be positive")
	}
	return out
}
func (p *Policy) AddLevel(level Sensitivity) {
	if !p.Allows(level) {
		p.Allowed = append(p.Allowed, level)
	}
}
func (p *Policy) RemoveLevel(level Sensitivity) {
	out := p.Allowed[:0]
	for _, v := range p.Allowed {
		if v != level {
			out = append(out, v)
		}
	}
	p.Allowed = out
}
