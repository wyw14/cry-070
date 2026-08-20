package masking

import "sort"

type Snapshot struct {
	ID, SourceID, Digest string
	Rows                 int
	Fields               []string
}

func (s Snapshot) Clone() Snapshot { s.Fields = append([]string(nil), s.Fields...); return s }
func (s *Snapshot) Normalize() {
	sort.Strings(s.Fields)
	if s.Rows < 0 {
		s.Rows = 0
	}
}
