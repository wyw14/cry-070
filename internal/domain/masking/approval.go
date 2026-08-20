package masking

import (
	"errors"
	"sort"
	"time"
)

var ErrApproval = errors.New("approval requirement not met")

type Approval struct {
	PipelineID, Actor, Decision, Note string
	At                                time.Time
}
type ApprovalLedger struct {
	Required []string
	Entries  []Approval
}

func (l *ApprovalLedger) Add(actor, decision, note string, now time.Time) error {
	if actor == "" || decision != "approve" && decision != "reject" {
		return ErrApproval
	}
	for _, e := range l.Entries {
		if e.Actor == actor {
			return ErrConflict
		}
	}
	l.Entries = append(l.Entries, Approval{Actor: actor, Decision: decision, Note: note, At: now.UTC()})
	return nil
}
func (l ApprovalLedger) Complete() bool {
	if len(l.Required) == 0 {
		return false
	}
	seen := map[string]bool{}
	for _, e := range l.Entries {
		if e.Decision != "approve" {
			return false
		}
		seen[e.Actor] = true
	}
	for _, r := range l.Required {
		if !seen[r] {
			return false
		}
	}
	return true
}
func (l ApprovalLedger) Pending() []string {
	seen := map[string]bool{}
	for _, e := range l.Entries {
		seen[e.Actor] = true
	}
	out := []string{}
	for _, r := range l.Required {
		if !seen[r] {
			out = append(out, r)
		}
	}
	sort.Strings(out)
	return out
}
