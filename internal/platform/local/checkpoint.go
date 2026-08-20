package local

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Checkpoint struct {
	BatchID string `json:"batch_id"`
	Offset  int    `json:"offset"`
	Digest  string `json:"digest"`
}

func SaveCheckpoint(root string, c Checkpoint) error {
	if err := os.MkdirAll(root, 0750); err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(root, c.BatchID+".json"), data, 0600)
}
func LoadCheckpoint(root, id string) (Checkpoint, error) {
	data, err := os.ReadFile(filepath.Join(root, id+".json"))
	if err != nil {
		return Checkpoint{}, err
	}
	var c Checkpoint
	err = json.Unmarshal(data, &c)
	return c, err
}
