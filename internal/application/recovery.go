package application

import (
	"context"
	"github.com/wyw14/cry052/internal/domain/masking"
)

type RecoveryService struct{}

func (RecoveryService) Resume(ctx context.Context, b masking.Batch, snapshot string) (masking.RecoveryDecision, error) {
	if err := ctx.Err(); err != nil {
		return masking.RecoveryDecision{}, err
	}
	d := masking.PlanRecovery(b, snapshot)
	if err := masking.ValidateRecovery(d, b); err != nil {
		return d, err
	}
	return d, nil
}
