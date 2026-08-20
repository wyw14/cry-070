package application

import (
	"context"
	"github.com/wyw14/cry052/internal/domain/masking"
	"time"
)

type AuditService struct {
	store interface {
		ListAudits(context.Context, string) []masking.AuditEntry
	}
	clock func() time.Time
}

func NewAuditService(store interface {
	ListAudits(context.Context, string) []masking.AuditEntry
}, clock func() time.Time) *AuditService {
	return &AuditService{store: store, clock: clock}
}
func (s *AuditService) Timeline(ctx context.Context, resource string) []masking.AuditEntry {
	if err := ctx.Err(); err != nil {
		return nil
	}
	out := s.store.ListAudits(ctx, resource)
	for i := range out {
		if out[i].CreatedAt.IsZero() {
			out[i].CreatedAt = s.clock()
		}
	}
	return out
}
