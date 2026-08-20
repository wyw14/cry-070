package tests

import (
	"github.com/wyw14/cry052/internal/domain/masking"
	"testing"
)

func TestSnapshotCloneDoesNotAliasFields(t *testing.T) {
	s := masking.Snapshot{Fields: []string{"email"}, Checksums: []string{"a"}, ColumnTypes: map[string]string{"email": "text"}}
	c := s.Clone()
	c.Fields[0] = "phone"
	c.Checksums[0] = "b"
	c.ColumnTypes["email"] = "binary"
	if s.Fields[0] != "email" {
		t.Fatal("snapshot clone aliases nested fields")
	}
	if s.Checksums[0] != "a" || s.ColumnTypes["email"] != "text" {
		t.Fatal("snapshot clone aliases nested metadata")
	}
}
