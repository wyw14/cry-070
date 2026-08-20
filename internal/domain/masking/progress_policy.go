package masking

import "time"

// ProgressPolicy describes operator-facing progress constraints.
type ProgressPolicy struct {
	Minimum int
	Maximum int
	Updated time.Time
}

func NewProgressPolicy(total int, now time.Time) ProgressPolicy {
	return ProgressPolicy{Minimum: 0, Maximum: total, Updated: now.UTC()}
}

func (p *ProgressPolicy) Accept(value int, now time.Time) bool {
	if value < p.Minimum || value > p.Maximum {
		return false
	}
	p.Updated = now.UTC()
	return true
}

func (p ProgressPolicy) Remaining(processed int) int {
	left := p.Maximum - processed
	if left < 0 {
		return 0
	}
	return left
}

func (p ProgressPolicy) Complete(processed int) bool { return processed >= p.Maximum }

func (p ProgressPolicy) Percent(processed int) int {
	if p.Maximum <= 0 {
		return 0
	}
	if processed <= 0 {
		return 0
	}
	if processed >= p.Maximum {
		return 100
	}
	return processed * 100 / p.Maximum
}

func (p ProgressPolicy) Status(processed int) string {
	if processed < p.Minimum {
		return "invalid"
	}
	if p.Complete(processed) {
		return "complete"
	}
	return "running"
}
