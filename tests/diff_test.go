package tests

import (
	"github.com/wyw14/cry052/internal/domain/masking"
	"testing"
)

func TestCompareRowsReportsUnknownTargetKey(t *testing.T) {
	got := masking.CompareRows([]map[string]string{{"id": "1", "email": "a"}}, []map[string]string{{"id": "2", "email": "b"}}, "id")
	if len(got) == 0 {
		t.Fatal("missing target row difference")
	}
}
