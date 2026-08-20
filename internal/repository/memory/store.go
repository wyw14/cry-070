package memory

import (
	"context"
	"github.com/wyw14/cry052/internal/domain/masking"
	"sort"
	"sync"
)

type Store struct {
	mu        sync.RWMutex
	sources   map[string]masking.Source
	fields    map[string]masking.Field
	rules     map[string]masking.Rule
	mappings  map[string]masking.Mapping
	pipelines map[string]masking.Pipeline
	previews  map[string]masking.Preview
	batches   map[string]masking.Batch
	audits    []masking.AuditEntry
	keys      map[string]string
}

func New() *Store {
	return &Store{sources: map[string]masking.Source{}, fields: map[string]masking.Field{}, rules: map[string]masking.Rule{}, mappings: map[string]masking.Mapping{}, pipelines: map[string]masking.Pipeline{}, previews: map[string]masking.Preview{}, batches: map[string]masking.Batch{}, keys: map[string]string{}}
}
func (s *Store) Within(ctx context.Context, fn func(*Store) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return fn(s)
}
func (s *Store) CreateSource(ctx context.Context, v masking.Source) error {
	return s.Within(ctx, func(_ *Store) error {
		if _, ok := s.sources[v.ID]; ok {
			return masking.ErrConflict
		}
		s.sources[v.ID] = v.Clone()
		return nil
	})
}
func (s *Store) GetSource(ctx context.Context, id string) (masking.Source, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.sources[id]
	if !ok {
		return masking.Source{}, masking.ErrInvalid
	}
	return v.Clone(), nil
}
func (s *Store) UpdateSource(ctx context.Context, v masking.Source, expected uint64) error {
	return s.Within(ctx, func(_ *Store) error {
		cur, ok := s.sources[v.ID]
		if !ok {
			return masking.ErrInvalid
		}
		if cur.Version != expected {
			return masking.ErrConflict
		}
		s.sources[v.ID] = v.Clone()
		return nil
	})
}
func (s *Store) ListSources(ctx context.Context, limit int) []masking.Source {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]masking.Source, 0, len(s.sources))
	for _, v := range s.sources {
		out = append(out, v.Clone())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}
func (s *Store) CreateField(ctx context.Context, v masking.Field) error {
	return s.Within(ctx, func(_ *Store) error {
		if _, ok := s.fields[v.ID]; ok {
			return masking.ErrConflict
		}
		s.fields[v.ID] = v.Clone()
		return nil
	})
}
func (s *Store) GetField(ctx context.Context, id string) (masking.Field, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.fields[id]
	if !ok {
		return masking.Field{}, masking.ErrInvalid
	}
	return v.Clone(), nil
}
func (s *Store) UpdateField(ctx context.Context, v masking.Field, expected uint64) error {
	return s.Within(ctx, func(_ *Store) error {
		cur, ok := s.fields[v.ID]
		if !ok {
			return masking.ErrInvalid
		}
		if cur.Version != expected {
			return masking.ErrConflict
		}
		s.fields[v.ID] = v.Clone()
		return nil
	})
}
func (s *Store) ListFields(ctx context.Context, source string) []masking.Field {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []masking.Field{}
	for _, v := range s.fields {
		if source == "" || v.SourceID == source {
			out = append(out, v.Clone())
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func (s *Store) CreateRule(ctx context.Context, v masking.Rule) error {
	return s.Within(ctx, func(_ *Store) error {
		if _, ok := s.rules[v.ID]; ok {
			return masking.ErrConflict
		}
		s.rules[v.ID] = v
		return nil
	})
}
func (s *Store) GetRule(ctx context.Context, id string) (masking.Rule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.rules[id]
	if !ok {
		return masking.Rule{}, masking.ErrInvalid
	}
	return v, nil
}
func (s *Store) ListRules(ctx context.Context) []masking.Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []masking.Rule{}
	for _, v := range s.rules {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func (s *Store) CreatePipeline(ctx context.Context, v masking.Pipeline) error {
	return s.Within(ctx, func(_ *Store) error {
		if _, ok := s.pipelines[v.ID]; ok {
			return masking.ErrConflict
		}
		s.pipelines[v.ID] = v.Clone()
		return nil
	})
}
func (s *Store) GetPipeline(ctx context.Context, id string) (masking.Pipeline, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.pipelines[id]
	if !ok {
		return masking.Pipeline{}, masking.ErrInvalid
	}
	return v.Clone(), nil
}
func (s *Store) UpdatePipeline(ctx context.Context, v masking.Pipeline, expected int) error {
	return s.Within(ctx, func(_ *Store) error {
		cur, ok := s.pipelines[v.ID]
		if !ok {
			return masking.ErrInvalid
		}
		if cur.Version != expected {
			return masking.ErrConflict
		}
		s.pipelines[v.ID] = v.Clone()
		return nil
	})
}
func (s *Store) CreateMapping(ctx context.Context, v masking.Mapping) error {
	return s.Within(ctx, func(_ *Store) error {
		if _, ok := s.mappings[v.ID]; ok {
			return masking.ErrConflict
		}
		s.mappings[v.ID] = v
		return nil
	})
}
func (s *Store) ListMappings(ctx context.Context, pipeline string) []masking.Mapping {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []masking.Mapping{}
	for _, v := range s.mappings {
		if pipeline == "" || v.PipelineID == pipeline {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Priority < out[j].Priority })
	return out
}
func (s *Store) CreatePreview(ctx context.Context, v masking.Preview) error {
	return s.Within(ctx, func(_ *Store) error {
		if _, ok := s.previews[v.ID]; ok {
			return masking.ErrConflict
		}
		s.previews[v.ID] = v
		return nil
	})
}
func (s *Store) GetPreview(ctx context.Context, id string) (masking.Preview, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.previews[id]
	if !ok {
		return masking.Preview{}, masking.ErrInvalid
	}
	return v, nil
}
func (s *Store) UpdatePreview(ctx context.Context, v masking.Preview) error {
	return s.Within(ctx, func(_ *Store) error { s.previews[v.ID] = v; return nil })
}
func (s *Store) CreateBatch(ctx context.Context, v masking.Batch) error {
	return s.Within(ctx, func(_ *Store) error {
		if old, ok := s.keys[v.IdempotencyKey]; ok && old == "" {
			return masking.ErrConflict
		}
		if _, ok := s.batches[v.ID]; ok {
			return masking.ErrConflict
		}
		s.keys[v.IdempotencyKey] = v.ID
		s.batches[v.ID] = v
		return nil
	})
}
func (s *Store) GetBatch(ctx context.Context, id string) (masking.Batch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.batches[id]
	if !ok {
		return masking.Batch{}, masking.ErrInvalid
	}
	return v, nil
}
func (s *Store) UpdateBatch(ctx context.Context, v masking.Batch) error {
	return s.Within(ctx, func(_ *Store) error { s.batches[v.ID] = v; return nil })
}
func (s *Store) AppendAudit(ctx context.Context, v masking.AuditEntry) error {
	return s.Within(ctx, func(_ *Store) error { s.audits = append(s.audits, v); return nil })
}
func (s *Store) ListAudits(ctx context.Context, resource string) []masking.AuditEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []masking.AuditEntry{}
	for _, v := range s.audits {
		if resource == "" || v.Resource == resource {
			out = append(out, v)
		}
	}
	return out
}
