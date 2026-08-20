package application

import (
	"context"
	"github.com/wyw14/cry052/internal/domain/masking"
)

type Planner struct{}

func (Planner) Plan(ctx context.Context, p masking.Pipeline, m []masking.Mapping, f []masking.Field, r []masking.Rule) (masking.ExecutionPlan, error) {
	if err := ctx.Err(); err != nil {
		return masking.ExecutionPlan{}, err
	}
	plan := masking.BuildPlan(p, m, f, r)
	return plan, nil
}
func (Planner) Ready(plan masking.ExecutionPlan) bool {
	return len(plan.Steps) > 0 && len(plan.Warnings) == 0
}
