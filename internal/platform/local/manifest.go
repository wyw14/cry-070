package local

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
)

type AuditManifest struct {
	Version int      `json:"version"`
	BatchID string   `json:"batch_id"`
	Files   []string `json:"files"`
	Digest  string   `json:"digest"`
}

func BuildAuditManifest(batchID string, files []string) AuditManifest {
	m := AuditManifest{Version: 1, BatchID: batchID, Files: append([]string(nil), files...)}
	data, _ := json.Marshal(m)
	sum := sha256.Sum256(data)
	m.Digest = hex.EncodeToString(sum[:])
	return m
}
func SaveAuditManifest(path string, m AuditManifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
func LoadAuditManifest(path string) (AuditManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return AuditManifest{}, err
	}
	var m AuditManifest
	err = json.Unmarshal(data, &m)
	return m, err
}
