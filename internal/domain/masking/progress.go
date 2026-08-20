package masking

import "time"

type Progress struct {
	BatchID          string
	Done, Total      int
	Started, Updated time.Time
	Message          string
}

func (p Progress) Fraction() float64 {
	if p.Total <= 0 {
		return 0
	}
	return float64(p.Done) / float64(p.Total)
}
func (p Progress) Complete() bool { return p.Total > 0 && p.Done >= p.Total }
func (p *Progress) Advance(n int, now time.Time) {
	if n < 0 {
		n = 0
	}
	p.Done += n
	if p.Done > p.Total {
		p.Done = p.Total
	}
	p.Updated = now.UTC()
}
