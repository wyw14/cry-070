package tests

import (
	"context"
	"github.com/wyw14/cry052/internal/application"
	"github.com/wyw14/cry052/internal/domain/masking"
	"github.com/wyw14/cry052/internal/repository/memory"
	"testing"
	"time"
)

func TestRunBatchHonorsCancelledContext(t *testing.T) {
	ctx := &cancelAfterFirstCheck{}
	svc := application.New(memory.New(), application.RealClock{}, &application.IDs{})
	b, _ := masking.NewBatch("b", "p", "v", "o", "t", "k", 1, time.Unix(1, 0))
	_ = svc.Store.CreateBatch(context.Background(), b)
	if _, err := svc.RunBatch(ctx, "b", "snap", 1); err == nil {
		t.Fatal("cancelled context ignored")
	}
}

type cancelAfterFirstCheck struct{ checks int }

func (c *cancelAfterFirstCheck) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c *cancelAfterFirstCheck) Done() <-chan struct{}       { return nil }
func (c *cancelAfterFirstCheck) Err() error {
	c.checks++
	if c.checks > 1 {
		return context.Canceled
	}
	return nil
}
func (c *cancelAfterFirstCheck) Value(any) any { return nil }
