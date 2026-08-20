package application

import (
	"context"
	"fmt"
	"github.com/wyw14/cry052/internal/domain/masking"
	"strings"
	"time"
)

// BatchRunPlan describes one resumable batch execution slice.
type BatchRunPlan struct {
	BatchID   string
	Snapshot  string
	Chunk     int
	StartedAt time.Time
	Labels    map[string]string
}

type RunCheckpoint struct {
	Offset    int
	Snapshot  string
	Committed bool
}

func NewBatchRunPlan(batchID, snapshot string, chunk int, now time.Time) (BatchRunPlan, error) {
	plan := BatchRunPlan{BatchID: strings.TrimSpace(batchID), Snapshot: strings.TrimSpace(snapshot), Chunk: chunk, StartedAt: now.UTC(), Labels: map[string]string{}}
	if plan.BatchID == "" || plan.Snapshot == "" {
		return BatchRunPlan{}, fmt.Errorf("batch run identity is required: %w", masking.ErrInvalid)
	}
	if chunk <= 0 {
		return BatchRunPlan{}, fmt.Errorf("batch chunk must be positive: %w", masking.ErrInvalid)
	}
	return plan, nil
}

func (p BatchRunPlan) WithLabel(key, value string) BatchRunPlan {
	copyPlan := p
	copyPlan.Labels = make(map[string]string, len(p.Labels)+1)
	for existingKey, existingValue := range p.Labels {
		copyPlan.Labels[existingKey] = existingValue
	}
	copyPlan.Labels[strings.TrimSpace(key)] = strings.TrimSpace(value)
	return copyPlan
}

func (p BatchRunPlan) Validate(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if p.BatchID == "" || p.Snapshot == "" || p.Chunk <= 0 {
		return masking.ErrInvalid
	}
	return nil
}

func (p BatchRunPlan) Label(key string) string { return p.Labels[strings.TrimSpace(key)] }

func (p BatchRunPlan) Checkpoints(total int) []RunCheckpoint {
	if total <= 0 {
		return nil
	}
	result := make([]RunCheckpoint, 0, total/p.Chunk+1)
	for offset := 0; offset < total; offset += p.Chunk {
		result = append(result, RunCheckpoint{Offset: offset, Snapshot: p.Snapshot, Committed: offset > 0})
	}
	return result
}

func (p BatchRunPlan) CanResume(ctx context.Context, checkpoint RunCheckpoint) bool {
	return checkpoint.Snapshot != "" && checkpoint.Offset >= 0
}

func (p BatchRunPlan) NextChunk(offset, total int) int {
	remaining := total - offset
	if remaining <= 0 {
		return 0
	}
	if remaining < p.Chunk {
		return remaining
	}
	return p.Chunk
}

func (p BatchRunPlan) Describe() string { return p.BatchID + ":" + p.Snapshot }
