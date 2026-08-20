package tests

import (
	"github.com/wyw14/cry052/internal/domain/masking"
	"testing"
	"time"
)

func TestSucceededBatchCannotBeCancelled(t *testing.T) {
	b, _ := masking.NewBatch("b", "p", "v", "o", "t", "k", 1, time.Unix(1, 0))
	_ = b.Start("snap", time.Unix(1, 0))
	_ = b.Advance(1, time.Unix(2, 0))
	if err := b.Cancel(time.Unix(3, 0)); err == nil {
		t.Fatal("succeeded batch cancelled")
	}
}
