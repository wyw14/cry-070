package application

import (
	"context"
	"github.com/wyw14/cry052/internal/domain/masking"
	"sort"
	"time"
)

type Report struct {
	Batch       masking.Batch
	Progress    masking.Progress
	Metrics     masking.Metrics
	GeneratedAt time.Time
	Notes       []string
}
type ReportBuilder struct{ Now func() time.Time }

func (b ReportBuilder) Build(ctx context.Context, batch masking.Batch, progress masking.Progress, metrics masking.Metrics) (Report, error) {
	if err := ctx.Err(); err != nil {
		return Report{}, err
	}
	now := time.Now().UTC()
	if b.Now != nil {
		now = b.Now()
	}
	r := Report{Batch: batch.Clone(), Progress: progress, Metrics: metrics, GeneratedAt: now, Notes: []string{}}
	if batch.State == masking.BatchFailed {
		r.Notes = append(r.Notes, "batch requires recovery")
	}
	if len(metrics.BySensitivity) > 0 {
		r.Notes = append(r.Notes, "sensitivity profile captured")
	}
	sort.Strings(r.Notes)
	return r, nil
}
func (r Report) ReadyForExport() bool {
	return r.Batch.State == masking.BatchSucceeded && r.Progress.Complete()
}
func (r Report) Headline() string {
	if r.ReadyForExport() {
		return "脱敏批次已完成"
	}
	if r.Batch.State == masking.BatchFailed {
		return "脱敏批次需要恢复"
	}
	return "脱敏批次处理中"
}
