package tests

import (
	"github.com/wyw14/cry052/internal/domain/masking"
	"testing"
)

func TestRunningReportIsNotSuccessful(t *testing.T) {
	if (masking.ExecutionReport{State: masking.BatchRunning}).Success() {
		t.Fatal("running report marked successful")
	}
}
