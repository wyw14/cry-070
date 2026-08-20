package tests

import (
	"context"
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
