package memory

import (
	"context"
	"github.com/wyw14/cry052/internal/domain/masking"
	"sort"
)

func (s *Store) FindFields(ctx context.Context, sourceID string, limit int) []masking.Field {
	if err := ctx.Err(); err != nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []masking.Field{}
	for _, f := range s.fields {
		if f.SourceID == sourceID {
			out = append(out, f.Clone())
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}
func (s *Store) FindRules(ctx context.Context, active bool) []masking.Rule {
	if err := ctx.Err(); err != nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []masking.Rule{}
	for _, r := range s.rules {
		if r.Active == active {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func (s *Store) FindBatches(ctx context.Context, state masking.BatchState) []masking.Batch {
	if err := ctx.Err(); err != nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []masking.Batch{}
	for _, b := range s.batches {
		if state == "" || b.State == state {
			out = append(out, b.Clone())
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartedAt.Before(out[j].StartedAt) })
	return out
}
