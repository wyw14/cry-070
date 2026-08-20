package masking

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var ErrInvalid = errors.New("invalid masking entity")
var ErrConflict = errors.New("masking conflict")
var ErrState = errors.New("invalid lifecycle state")

type SourceState string

const (
	SourceDraft    SourceState = "draft"
	SourceReady    SourceState = "ready"
	SourceDisabled SourceState = "disabled"
)

type Sensitivity string

const (
	Public     Sensitivity = "public"
	Internal   Sensitivity = "internal"
	Sensitive  Sensitivity = "sensitive"
	Restricted Sensitivity = "restricted"
)

type RuleKind string

const (
	Mask       RuleKind = "mask"
	Replace    RuleKind = "replace"
	Generalize RuleKind = "generalize"
	Hash       RuleKind = "hash"
	Preserve   RuleKind = "preserve"
)

type PreviewState string

const (
	PreviewPending   PreviewState = "pending"
	PreviewConfirmed PreviewState = "confirmed"
	PreviewRejected  PreviewState = "rejected"
)

type BatchState string

const (
	BatchPlanned    BatchState = "planned"
	BatchRunning    BatchState = "running"
	BatchPaused     BatchState = "paused"
	BatchSucceeded  BatchState = "succeeded"
	BatchFailed     BatchState = "failed"
	BatchRolledBack BatchState = "rolled_back"
	BatchCancelled  BatchState = "cancelled"
)

type Source struct {
	ID, Name, Adapter, Owner string
	EncryptedDSN             string
	State                    SourceState
	Tables                   []string
	CreatedAt, UpdatedAt     time.Time
	Version                  uint64
}

func NewSource(id, name, adapter, owner string, now time.Time) (Source, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" || strings.TrimSpace(adapter) == "" || strings.TrimSpace(owner) == "" {
		return Source{}, ErrInvalid
	}
	return Source{ID: id, Name: name, Adapter: adapter, Owner: owner, State: SourceDraft, CreatedAt: now.UTC(), UpdatedAt: now.UTC(), Version: 1}, nil
}
func (s *Source) RegisterTable(name string, now time.Time) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalid
	}
	for _, v := range s.Tables {
		if v == name {
			return ErrConflict
		}
	}
	s.Tables = append(s.Tables, name)
	sort.Strings(s.Tables)
	s.touch(now)
	return nil
}
func (s *Source) EncryptPlaceholder(value string, now time.Time) error {
	if strings.TrimSpace(value) == "" {
		return ErrInvalid
	}
	sum := sha256.Sum256([]byte(value))
	s.EncryptedDSN = hex.EncodeToString(sum[:])
	s.State = SourceReady
	s.touch(now)
	return nil
}
func (s *Source) Disable(now time.Time) { s.State = SourceDisabled; s.touch(now) }
func (s *Source) touch(now time.Time)   { s.UpdatedAt = now.UTC(); s.Version++ }
func (s Source) Clone() Source          { s.Tables = append([]string(nil), s.Tables...); return s }

type Field struct {
	ID, SourceID, Table, Name string
	Type                      string
	Sensitivity               Sensitivity
	AccessScope               []string
	Nullable                  bool
	UpdatedAt                 time.Time
	Version                   uint64
}

func NewField(id, source, table, name, typ string, level Sensitivity, now time.Time) (Field, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(source) == "" || strings.TrimSpace(table) == "" || strings.TrimSpace(name) == "" {
		return Field{}, ErrInvalid
	}
	if typ == "" {
		typ = "text"
	}
	return Field{ID: id, SourceID: source, Table: table, Name: name, Type: typ, Sensitivity: level, UpdatedAt: now.UTC(), Version: 1}, nil
}
func (f *Field) SetScope(scope []string, now time.Time) error {
	if len(scope) == 0 && f.Sensitivity == Restricted {
		return fmt.Errorf("restricted field needs scope: %w", ErrInvalid)
	}
	f.AccessScope = unique(scope)
	f.UpdatedAt = now.UTC()
	f.Version++
	return nil
}
func (f Field) Clone() Field { f.AccessScope = append([]string(nil), f.AccessScope...); return f }

type Rule struct {
	ID, Name, Description string
	Kind                  RuleKind
	Pattern, Replacement  string
	Version               int
	Active                bool
	CreatedBy             string
	CreatedAt             time.Time
}

func NewRule(id, name string, kind RuleKind, pattern, replacement, owner string, now time.Time) (Rule, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" || owner == "" {
		return Rule{}, ErrInvalid
	}
	if !validKind(kind) {
		return Rule{}, ErrInvalid
	}
	if kind == Mask && pattern == "" {
		pattern = "first"
	}
	return Rule{ID: id, Name: name, Kind: kind, Pattern: pattern, Replacement: replacement, Version: 1, Active: true, CreatedBy: owner, CreatedAt: now.UTC()}, nil
}
func (r *Rule) Validate(value string) (string, error) {
	if !r.Active {
		return value, ErrState
	}
	switch r.Kind {
	case Mask:
		return maskValue(value, r.Pattern), nil
	case Replace:
		if r.Replacement == "" {
			return "", ErrInvalid
		}
		return r.Replacement, nil
	case Generalize:
		return generalize(value, r.Pattern), nil
	case Hash:
		sum := sha256.Sum256([]byte(value))
		return hex.EncodeToString(sum[:]), nil
	case Preserve:
		return value, nil
	}
	return "", ErrInvalid
}
func (r *Rule) Deactivate() { r.Active = false; r.Version++ }
func validKind(k RuleKind) bool {
	return k == Mask || k == Replace || k == Generalize || k == Hash || k == Preserve
}
func maskValue(v, mode string) string {
	if v == "" {
		return v
	}
	if mode == "all" {
		return strings.Repeat("*", len([]rune(v)))
	}
	r := []rune(v)
	if len(r) <= 2 {
		return strings.Repeat("*", len(r))
	}
	return string(r[:1]) + strings.Repeat("*", len(r)-2) + string(r[len(r)-1])
}
func generalize(v, mode string) string {
	if mode == "year" && len(v) >= 4 {
		return v[:4]
	}
	if mode == "domain" {
		if i := strings.Index(v, "@"); i > 0 {
			return "*@" + v[i+1:]
		}
	}
	return "[generalized]"
}

type Mapping struct {
	ID, PipelineID, FieldID, RuleID string
	Priority                        int
	Required                        bool
	CreatedAt                       time.Time
}

func NewMapping(id, pipeline, field, rule string, priority int, now time.Time) (Mapping, error) {
	if id == "" || pipeline == "" || field == "" || rule == "" || priority < 0 {
		return Mapping{}, ErrInvalid
	}
	return Mapping{ID: id, PipelineID: pipeline, FieldID: field, RuleID: rule, Priority: priority, CreatedAt: now.UTC()}, nil
}

type Pipeline struct {
	ID, Name, Owner      string
	SourceID             string
	MappingIDs           []string
	Version              int
	State                SourceState
	ApprovedBy           string
	ApprovedAt           time.Time
	CreatedAt, UpdatedAt time.Time
}

func NewPipeline(id, name, owner, source string, now time.Time) (Pipeline, error) {
	if id == "" || name == "" || owner == "" || source == "" {
		return Pipeline{}, ErrInvalid
	}
	return Pipeline{ID: id, Name: name, Owner: owner, SourceID: source, State: SourceDraft, Version: 1, CreatedAt: now.UTC(), UpdatedAt: now.UTC()}, nil
}
func (p *Pipeline) AddMapping(id string, now time.Time) error {
	if id == "" {
		return ErrInvalid
	}
	for _, v := range p.MappingIDs {
		if v == id {
			return ErrConflict
		}
	}
	p.MappingIDs = append(p.MappingIDs, id)
	p.Version++
	p.UpdatedAt = now.UTC()
	return nil
}
func (p *Pipeline) Approve(actor string, now time.Time) error {
	if p.State != SourceDraft && p.State != SourceReady {
		return ErrState
	}
	if actor == "" || len(p.MappingIDs) == 0 {
		return ErrInvalid
	}
	p.State = SourceReady
	p.ApprovedBy = actor
	p.ApprovedAt = now.UTC()
	p.Version++
	p.UpdatedAt = now.UTC()
	return nil
}
func (p *Pipeline) Revoke(now time.Time) { p.State = SourceDraft; p.Version++; p.UpdatedAt = now.UTC() }
func (p Pipeline) Clone() Pipeline       { p.MappingIDs = append([]string(nil), p.MappingIDs...); return p }

type Preview struct {
	ID, PipelineID, SourceID string
	RowsIn, RowsOut          int
	Sample                   []map[string]string
	Conflicts                []string
	State                    PreviewState
	ConfirmedBy              string
	CreatedAt, ConfirmedAt   time.Time
}

func NewPreview(id, pipeline, source string, rows int, now time.Time) (Preview, error) {
	if id == "" || pipeline == "" || source == "" || rows < 0 {
		return Preview{}, ErrInvalid
	}
	return Preview{ID: id, PipelineID: pipeline, SourceID: source, RowsIn: rows, State: PreviewPending, CreatedAt: now.UTC()}, nil
}
func (p *Preview) Confirm(actor string, now time.Time) error {
	if p.State != PreviewPending || actor == "" || len(p.Conflicts) > 0 {
		return ErrState
	}
	p.State = PreviewConfirmed
	p.ConfirmedBy = actor
	p.ConfirmedAt = now.UTC()
	return nil
}
func (p *Preview) Reject(reason string) error {
	if p.State != PreviewPending {
		return ErrState
	}
	if reason != "" {
		p.Conflicts = append(p.Conflicts, reason)
	}
	p.State = PreviewRejected
	return nil
}

type Batch struct {
	ID, PipelineID, PreviewID, Owner string
	State                            BatchState
	Total, Processed, Failed         int
	InputSnapshot, Target            string
	IdempotencyKey                   string
	ErrorSummary                     string
	StartedAt, FinishedAt            time.Time
	Version                          uint64
}

func NewBatch(id, pipeline, preview, owner, target, key string, total int, now time.Time) (Batch, error) {
	if id == "" || pipeline == "" || preview == "" || owner == "" || target == "" || key == "" || total < 0 {
		return Batch{}, ErrInvalid
	}
	return Batch{ID: id, PipelineID: pipeline, PreviewID: preview, Owner: owner, Target: target, IdempotencyKey: key, Total: total, State: BatchPlanned, Version: 1, StartedAt: now.UTC()}, nil
}
func (b *Batch) Start(snapshot string, now time.Time) error {
	if b.State != BatchPlanned || snapshot == "" {
		return ErrState
	}
	b.InputSnapshot = snapshot
	b.State = BatchRunning
	b.Version++
	b.StartedAt = now.UTC()
	return nil
}
func (b *Batch) Advance(count int, now time.Time) error {
	if b.State != BatchRunning || count < 0 {
		return ErrState
	}
	b.Processed += count
	if b.Processed > b.Total {
		b.Processed = b.Total
	}
	b.Version++
	if b.Processed == b.Total {
		b.State = BatchSucceeded
		b.FinishedAt = now.UTC()
	}
	return nil
}
func (b *Batch) Fail(reason string, now time.Time) {
	b.State = BatchFailed
	b.ErrorSummary = strings.TrimSpace(reason)
	b.FinishedAt = now.UTC()
	b.Version++
}
// Cancel transitions a batch to the cancelled state. Only batches that are
// still in flight (planned or running) may be cancelled; a succeeded batch is
// terminal and must not be cancellable. This mirrors the guard in
// CancellationPolicy.CanCancel so the domain and the API-facing policy agree.
func (b *Batch) Cancel(now time.Time) error {
	if b.State != BatchRunning && b.State != BatchPlanned {
		return ErrState
	}
	b.State = BatchCancelled
	b.FinishedAt = now.UTC()
	b.Version++
	return nil
}
func (b *Batch) Rollback(now time.Time) error {
	if b.State != BatchFailed && b.State != BatchSucceeded {
		return ErrState
	}
	b.State = BatchRolledBack
	b.FinishedAt = now.UTC()
	b.Version++
	return nil
}
func (b Batch) Clone() Batch { return b }

type AuditEntry struct {
	ID, Actor, Action, Resource, ResourceID, RequestID string
	Redacted                                           bool
	CreatedAt                                          time.Time
}

func NewAudit(id, actor, action, resource, resourceID, requestID string, now time.Time) (AuditEntry, error) {
	if id == "" || actor == "" || action == "" || resource == "" {
		return AuditEntry{}, ErrInvalid
	}
	return AuditEntry{ID: id, Actor: actor, Action: action, Resource: resource, ResourceID: resourceID, RequestID: requestID, CreatedAt: now.UTC()}, nil
}
func (e *AuditEntry) Redact() { e.Actor = "[redacted]"; e.ResourceID = "[redacted]"; e.Redacted = true }

func unique(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}

var identifierPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]{0,63}$`)

func ValidateIdentifier(value string) error {
	if !identifierPattern.MatchString(value) {
		return fmt.Errorf("invalid identifier: %w", ErrInvalid)
	}
	return nil
}
