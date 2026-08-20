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

func TestCompareRowsReportsAddedRemovedAndChanged(t *testing.T) {
	before := []map[string]string{
		{"id": "1", "email": "a"},
		{"id": "2", "email": "b"},
	}
	after := []map[string]string{
		{"id": "1", "email": "changed"},
		{"id": "3", "email": "c"},
	}
	got := masking.CompareRows(before, after, "id")
	if len(got) != 3 {
		t.Fatalf("expected 3 row diffs, got %d: %+v", len(got), got)
	}
	byKey := map[string]masking.RowDiff{}
	for _, d := range got {
		byKey[d.Key] = d
	}
	if d := byKey["1"]; d.Status != "changed" || len(d.Changes) != 1 || d.Changes[0].After != "changed" {
		t.Fatalf("id=1 not reported as changed: %+v", d)
	}
	if d := byKey["2"]; d.Status != "removed" {
		t.Fatalf("id=2 not reported as removed: %+v", d)
	}
	if d := byKey["3"]; d.Status != "added" || len(d.Changes) != 1 || d.Changes[0].After != "c" {
		t.Fatalf("id=3 not reported as added: %+v", d)
	}
}
