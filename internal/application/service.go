package application

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/wyw14/cry052/internal/domain/masking"
	"github.com/wyw14/cry052/internal/repository/memory"
	"strings"
	"sync/atomic"
	"time"
)

type Clock interface{ Now() time.Time }
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now().UTC() }

type IDs struct{ n atomic.Uint64 }

func (i *IDs) New(prefix string) string { return fmt.Sprintf("%s-%d", prefix, i.n.Add(1)) }

type Services struct {
	Store    *memory.Store
	Clock    Clock
	IDs      *IDs
	Validate *validator.Validate
}

func New(store *memory.Store, clock Clock, ids *IDs) *Services {
	if clock == nil {
		clock = RealClock{}
	}
	if ids == nil {
		ids = &IDs{}
	}
	return &Services{Store: store, Clock: clock, IDs: ids, Validate: validator.New()}
}
func (s *Services) RegisterSource(ctx context.Context, name, adapter, owner, dsn string) (masking.Source, error) {
	v, err := masking.NewSource(s.IDs.New("src"), name, adapter, owner, s.Clock.Now())
	if err != nil {
		return v, err
	}
	if err = s.Store.CreateSource(ctx, v); err != nil {
		return v, err
	}
	if dsn != "" {
		v, _ = s.Store.GetSource(ctx, v.ID)
		if err = v.EncryptPlaceholder(dsn, s.Clock.Now()); err != nil {
			return v, err
		}
		err = s.Store.UpdateSource(ctx, v, v.Version-1)
	}
	return v, err
}
func (s *Services) RegisterTable(ctx context.Context, source, table string) error {
	v, err := s.Store.GetSource(ctx, source)
	if err != nil {
		return err
	}
	before := v.Version
	if err = v.RegisterTable(table, s.Clock.Now()); err != nil {
		return err
	}
	return s.Store.UpdateSource(ctx, v, before)
}
func (s *Services) AddField(ctx context.Context, source, table, name, typ string, level masking.Sensitivity) (masking.Field, error) {
	if err := masking.ValidateIdentifier(table); err != nil {
		return masking.Field{}, err
	}
	v, err := masking.NewField(s.IDs.New("field"), source, table, name, typ, level, s.Clock.Now())
	if err != nil {
		return v, err
	}
	return v, s.Store.CreateField(ctx, v)
}
func (s *Services) ScopeField(ctx context.Context, id string, scope []string) error {
	v, err := s.Store.GetField(ctx, id)
	if err != nil {
		return err
	}
	before := v.Version
	if err = v.SetScope(scope, s.Clock.Now()); err != nil {
		return err
	}
	return s.Store.UpdateField(ctx, v, before)
}
func (s *Services) CreateRule(ctx context.Context, name string, kind masking.RuleKind, pattern, replacement, owner string) (masking.Rule, error) {
	v, err := masking.NewRule(s.IDs.New("rule"), name, kind, pattern, replacement, owner, s.Clock.Now())
	if err != nil {
		return v, err
	}
	return v, s.Store.CreateRule(ctx, v)
}
func (s *Services) CreatePipeline(ctx context.Context, name, owner, source string) (masking.Pipeline, error) {
	v, err := masking.NewPipeline(s.IDs.New("pipe"), name, owner, source, s.Clock.Now())
	if err != nil {
		return v, err
	}
	return v, s.Store.CreatePipeline(ctx, v)
}
func (s *Services) MapRule(ctx context.Context, pipeline, field, rule string, priority int) error {
	v, err := masking.NewMapping(s.IDs.New("map"), pipeline, field, rule, priority, s.Clock.Now())
	if err != nil {
		return err
	}
	if err = s.Store.CreateMapping(ctx, v); err != nil {
		return err
	}
	p, err := s.Store.GetPipeline(ctx, pipeline)
	if err != nil {
		return err
	}
	before := p.Version
	if err = p.AddMapping(v.ID, s.Clock.Now()); err != nil {
		return err
	}
	return s.Store.UpdatePipeline(ctx, p, before)
}
func (s *Services) ApprovePipeline(ctx context.Context, pipeline, actor string) error {
	p, err := s.Store.GetPipeline(ctx, pipeline)
	if err != nil {
		return err
	}
	before := p.Version
	if err = p.Approve(actor, s.Clock.Now()); err != nil {
		return err
	}
	return s.Store.UpdatePipeline(ctx, p, before)
}
func (s *Services) Preview(ctx context.Context, pipeline, source string, rows int) (masking.Preview, error) {
	p, err := s.Store.GetPipeline(ctx, pipeline)
	if err != nil {
		return masking.Preview{}, err
	}
	if p.State != masking.SourceReady {
		return masking.Preview{}, masking.ErrState
	}
	v, err := masking.NewPreview(s.IDs.New("preview"), p.ID, source, rows, s.Clock.Now())
	if err != nil {
		return v, err
	}
	fields := s.Store.ListFields(ctx, source)
	rules := s.Store.ListRules(ctx)
	for _, f := range fields {
		if f.Sensitivity == masking.Restricted && len(f.AccessScope) == 0 {
			v.Conflicts = append(v.Conflicts, "restricted field without scope: "+f.Name)
		}
	}
	if len(rules) == 0 {
		v.Conflicts = append(v.Conflicts, "pipeline has no active rule")
	}
	v.RowsOut = rows
	return v, s.Store.CreatePreview(ctx, v)
}
func (s *Services) ConfirmPreview(ctx context.Context, id, actor string) error {
	v, err := s.Store.GetPreview(ctx, id)
	if err != nil {
		return err
	}
	if err = v.Confirm(actor, s.Clock.Now()); err != nil {
		return err
	}
	return s.Store.UpdatePreview(ctx, v)
}
func (s *Services) StartBatch(ctx context.Context, pipeline, preview, owner, target, key string, total int) (masking.Batch, error) {
	p, err := s.Store.GetPipeline(ctx, pipeline)
	if err != nil {
		return masking.Batch{}, err
	}
	if p.State != masking.SourceReady {
		return masking.Batch{}, masking.ErrState
	}
	v, err := s.Store.GetPreview(ctx, preview)
	if err != nil {
		return masking.Batch{}, err
	}
	if v.State != masking.PreviewConfirmed {
		return masking.Batch{}, masking.ErrState
	}
	b, err := masking.NewBatch(s.IDs.New("batch"), pipeline, preview, owner, target, key, total, s.Clock.Now())
	if err != nil {
		return b, err
	}
	if err = s.Store.CreateBatch(ctx, b); err != nil {
		return b, err
	}
	audit, _ := masking.NewAudit(s.IDs.New("audit"), owner, "batch.created", "batch", b.ID, "", s.Clock.Now())
	_ = s.Store.AppendAudit(ctx, audit)
	return b, nil
}
func (s *Services) RunBatch(ctx context.Context, id, snapshot string, chunk int) (masking.Batch, error) {
	// BUG: cancellation is checked only after the state transition.
	b, err := s.Store.GetBatch(ctx, id)
	if err != nil {
		return b, err
	}
	plan, err := NewBatchRunPlan(id, snapshot, chunk, s.Clock.Now())
	if err != nil {
		return b, err
	}
	if err = b.Start(plan.Snapshot, plan.StartedAt); err != nil {
		return b, err
	}
	if err = s.Store.UpdateBatch(ctx, b); err != nil {
		return b, err
	}
	if err = b.Advance(plan.Chunk, s.Clock.Now()); err != nil {
		return b, err
	}
	_ = s.Store.UpdateBatch(ctx, b)
	return b, nil
}
func (s *Services) CancelBatch(ctx context.Context, id string) error {
	b, err := s.Store.GetBatch(ctx, id)
	if err != nil {
		return err
	}
	if err = b.Cancel(s.Clock.Now()); err != nil {
		return err
	}
	return s.Store.UpdateBatch(ctx, b)
}
func (s *Services) ApplyRule(ctx context.Context, ruleID, value string) (string, error) {
	r, err := s.Store.GetRule(ctx, ruleID)
	if err != nil {
		return "", err
	}
	return r.Validate(strings.TrimSpace(value))
}
func (s *Services) ExportAudit(ctx context.Context, resource string) []masking.AuditEntry {
	return s.Store.ListAudits(ctx, resource)
}
