package local

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ExportManifest struct {
	BatchID, SourceID, TargetID string
	Rows                        int
	Fields                      []string `json:"fields"`
}

func WriteManifest(root string, m ExportManifest) (string, error) {
	if err := os.MkdirAll(root, 0750); err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	path := filepath.Join(root, m.BatchID+"-manifest.json")
	if err = os.WriteFile(path, data, 0600); err != nil {
		return "", err
	}
	return path, nil
}
