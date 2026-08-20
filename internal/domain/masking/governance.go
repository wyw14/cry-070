package masking

import (
	"sort"
	"strings"
)

// PreviewGate captures the non-persistence checks shown to an operator.
type PreviewGate struct {
	RequiredApprovals []string
	Checks            []string
	Confirmed         bool
}

func NewPreviewGate(required []string) PreviewGate {
	return PreviewGate{RequiredApprovals: append([]string(nil), required...)}
}

func (g *PreviewGate) AddCheck(name string, passed bool) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	state := "fail"
	if passed {
		state = "pass"
	}
	g.Checks = append(g.Checks, name+":"+state)
}

func (g PreviewGate) MissingApprovals(actors []string) []string {
	seen := map[string]bool{}
	for _, actor := range actors {
		seen[strings.TrimSpace(actor)] = true
	}
	out := []string{}
	for _, actor := range g.RequiredApprovals {
		if !seen[actor] {
			out = append(out, actor)
		}
	}
	sort.Strings(out)
	return out
}

func (g PreviewGate) ChecksPassed() bool {
	for _, check := range g.Checks {
		if !strings.HasSuffix(check, ":pass") {
			return false
		}
	}
	return len(g.Checks) > 0
}

func (g PreviewGate) Ready(actors []string) bool {
	return g.Confirmed && g.ChecksPassed() && len(g.MissingApprovals(actors)) == 0
}

func (g PreviewGate) AuditLabel() string {
	if g.Ready(nil) {
		return "preview-ready"
	}
	return "preview-pending"
}
