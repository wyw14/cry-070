package tests

import (
	"context"
	"github.com/wyw14/cry052/internal/application"
	"github.com/wyw14/cry052/internal/domain/masking"
	"github.com/wyw14/cry052/internal/repository/memory"
	"testing"
	"time"
)

func TestBatchIdempotencyRejectsDifferentBatch(t *testing.T) {
	s := memory.New()
	a, _ := masking.NewBatch("a", "p", "v", "o", "t", "same", 1, timeNow())
	b, _ := masking.NewBatch("b", "p", "v", "o", "t", "same", 1, timeNow())
	if err := s.CreateBatch(context.Background(), a); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateBatch(context.Background(), b); err == nil {
		t.Fatal("duplicate idempotency key accepted")
	}
}

func timeNow() time.Time { return time.Unix(1, 0).UTC() }

// seedForBatch wires up a fully approved pipeline + confirmed preview so that
// StartBatch is allowed to proceed. It returns the values needed to call
// StartBatch.
func seedForBatch(t *testing.T, ctx context.Context, svc *application.Services) (pipeline, preview, owner, target string) {
	t.Helper()
	src, err := svc.RegisterSource(ctx, "样例库", "sqlite", "admin", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if err = svc.RegisterTable(ctx, src.ID, "customers"); err != nil {
		t.Fatal(err)
	}
	field, err := svc.AddField(ctx, src.ID, "customers", "email", "text", masking.Sensitive)
	if err != nil {
		t.Fatal(err)
	}
	rule, err := svc.CreateRule(ctx, "邮箱哈希", masking.Hash, "", "", "admin")
	if err != nil {
		t.Fatal(err)
	}
	pipe, err := svc.CreatePipeline(ctx, "客户脱敏", "admin", src.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err = svc.MapRule(ctx, pipe.ID, field.ID, rule.ID, 1); err != nil {
		t.Fatal(err)
	}
	if err = svc.ApprovePipeline(ctx, pipe.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	pv, err := svc.Preview(ctx, pipe.ID, src.ID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if err = svc.ConfirmPreview(ctx, pv.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	return pipe.ID, pv.ID, "admin", "target"
}

func TestStartBatchRetryReturnsSameBatch(t *testing.T) {
	ctx := context.Background()
	svc := application.New(memory.New(), application.RealClock{}, &application.IDs{})
	pipe, preview, owner, target := seedForBatch(t, ctx, svc)

	first, err := svc.StartBatch(ctx, pipe, preview, owner, target, "idem-retry", 2)
	if err != nil {
		t.Fatalf("first StartBatch failed: %v", err)
	}
	// A network retry submits the same idempotency key again.
	retry, err := svc.StartBatch(ctx, pipe, preview, owner, target, "idem-retry", 2)
	if err != nil {
		t.Fatalf("retry StartBatch failed: %v", err)
	}
	if retry.ID != first.ID {
		t.Fatalf("retry created a different batch: %q != %q", retry.ID, first.ID)
	}
}

func TestStartBatchDistinctKeysCreateDistinctBatches(t *testing.T) {
	ctx := context.Background()
	svc := application.New(memory.New(), application.RealClock{}, &application.IDs{})
	pipe, preview, owner, target := seedForBatch(t, ctx, svc)

	a, err := svc.StartBatch(ctx, pipe, preview, owner, target, "key-a", 2)
	if err != nil {
		t.Fatal(err)
	}
	b, err := svc.StartBatch(ctx, pipe, preview, owner, target, "key-b", 2)
	if err != nil {
		t.Fatal(err)
	}
	if a.ID == b.ID {
		t.Fatalf("distinct keys produced the same batch: %q", a.ID)
	}
}
