package local

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SnapshotStore struct{ root string }

func NewSnapshotStore(root string) (*SnapshotStore, error) {
	if root == "" {
		return nil, fmt.Errorf("snapshot root required")
	}
	if err := os.MkdirAll(root, 0750); err != nil {
		return nil, err
	}
	return &SnapshotStore{root: root}, nil
}
func (s *SnapshotStore) Save(id string, rows []map[string]string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("snapshot id required")
	}
	path := filepath.Join(s.root, id+".snapshot")
	if filepath.Dir(path) != filepath.Clean(s.root) {
		return "", fmt.Errorf("invalid snapshot path")
	}
	var b strings.Builder
	for _, row := range rows {
		for k, v := range row {
			fmt.Fprintf(&b, "%s=%s\n", k, v)
		}
	}
	sum := sha256.Sum256([]byte(b.String()))
	if err := os.WriteFile(path, []byte(b.String()), 0640); err != nil {
		return "", err
	}
	return hex.EncodeToString(sum[:]), nil
}
