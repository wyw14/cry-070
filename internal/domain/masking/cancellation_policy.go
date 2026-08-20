package masking

import "time"

// CancellationPolicy exposes cancellation decisions to API clients.
type CancellationPolicy struct {
	RequestedBy string
	Reason      string
	At          time.Time
}

func NewCancellationPolicy(actor, reason string, now time.Time) CancellationPolicy {
	return CancellationPolicy{RequestedBy: actor, Reason: reason, At: now.UTC()}
}

func (p CancellationPolicy) Valid() bool { return p.RequestedBy != "" && p.Reason != "" }

func (p CancellationPolicy) CanCancel(state BatchState) bool {
	return state == BatchPlanned || state == BatchRunning
}

func (p CancellationPolicy) Label(state BatchState) string {
	if p.CanCancel(state) {
		return "cancellable"
	}
	return "terminal"
}

func (p CancellationPolicy) AuditFields() []string {
	return []string{p.RequestedBy, p.Reason, p.At.Format(time.RFC3339)}
}
