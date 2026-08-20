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
	if b.State != masking.BatchSucceeded {
		t.Fatalf("precondition: state %s, want succeeded", b.State)
	}
	if err := b.Cancel(time.Unix(3, 0)); err == nil {
		t.Fatal("succeeded batch cancelled")
	}
	if b.State != masking.BatchSucceeded {
		t.Fatalf("succeeded state mutated to %s", b.State)
	}
}

func TestPlannedAndRunningBatchesRemainCancellable(t *testing.T) {
	planned, _ := masking.NewBatch("b1", "p", "v", "o", "t", "k", 1, time.Unix(1, 0))
	if err := planned.Cancel(time.Unix(2, 0)); err != nil {
		t.Fatalf("planned batch not cancellable: %v", err)
	}
	if planned.State != masking.BatchCancelled {
		t.Fatalf("planned state %s, want cancelled", planned.State)
	}

	running, _ := masking.NewBatch("b2", "p", "v", "o", "t", "k", 2, time.Unix(1, 0))
	if err := running.Start("snap", time.Unix(2, 0)); err != nil {
		t.Fatal(err)
	}
	if err := running.Cancel(time.Unix(3, 0)); err != nil {
		t.Fatalf("running batch not cancellable: %v", err)
	}
	if running.State != masking.BatchCancelled {
		t.Fatalf("running state %s, want cancelled", running.State)
	}
}
