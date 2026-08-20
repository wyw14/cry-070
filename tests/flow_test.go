package tests

import (
	"context"
	"github.com/wyw14/cry052/internal/application"
	"github.com/wyw14/cry052/internal/domain/masking"
	"github.com/wyw14/cry052/internal/repository/memory"
	"testing"
)

func TestPreviewMustBeConfirmedBeforeBatch(t *testing.T) {
	ctx := context.Background()
	svc := application.New(memory.New(), application.RealClock{}, &application.IDs{})
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
	if err = svc.ScopeField(ctx, field.ID, []string{"masking-worker"}); err != nil {
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
	preview, err := svc.Preview(ctx, pipe.ID, src.ID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.StartBatch(ctx, pipe.ID, preview.ID, "admin", "target", "idem-1", 2); err == nil {
		t.Fatal("unconfirmed preview accepted")
	}
	if err = svc.ConfirmPreview(ctx, preview.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	batch, err := svc.StartBatch(ctx, pipe.ID, preview.ID, "admin", "target", "idem-1", 2)
	if err != nil {
		t.Fatal(err)
	}
	if batch.State != masking.BatchPlanned {
		t.Fatalf("unexpected state %s", batch.State)
	}
}
func TestRuleMasksAndHashesWithoutExternalIO(t *testing.T) {
	ctx := context.Background()
	svc := application.New(memory.New(), application.RealClock{}, &application.IDs{})
	r, err := svc.CreateRule(ctx, "邮箱哈希", masking.Hash, "", "", "admin")
	if err != nil {
		t.Fatal(err)
	}
	out, err := svc.ApplyRule(ctx, r.ID, "alice@example.test")
	if err != nil || out == "alice@example.test" {
		t.Fatalf("hash rule failed: %q %v", out, err)
	}
}
