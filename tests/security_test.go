package tests

import (
	"context"
	"testing"

	"github.com/wyw14/cry052/internal/application"
	"github.com/wyw14/cry052/internal/domain/masking"
	"github.com/wyw14/cry052/internal/repository/memory"
)

func TestRestrictedFieldRequiresExplicitScope(t *testing.T) {
	ctx := context.Background()
	svc := application.New(memory.New(), application.RealClock{}, &application.IDs{})
	src, err := svc.RegisterSource(ctx, "audit-source", "sqlite", "admin", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.RegisterTable(ctx, src.ID, "accounts"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AddField(ctx, src.ID, "accounts", "phone", "text", masking.Restricted); err != nil {
		t.Fatal(err)
	}
	rule, err := svc.CreateRule(ctx, "mask-phone", masking.Mask, "", "***", "admin")
	if err != nil {
		t.Fatal(err)
	}
	pipeline, err := svc.CreatePipeline(ctx, "restricted-preview", "admin", src.ID)
	if err != nil {
		t.Fatal(err)
	}
	field, err := svc.AddField(ctx, src.ID, "accounts", "email", "text", masking.Sensitive)
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.MapRule(ctx, pipeline.ID, field.ID, rule.ID, 1); err != nil {
		t.Fatal(err)
	}
	if err := svc.ApprovePipeline(ctx, pipeline.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	preview, err := svc.Preview(ctx, pipeline.ID, src.ID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.Conflicts) == 0 {
		t.Fatal("restricted field without scope was not reported")
	}
}
