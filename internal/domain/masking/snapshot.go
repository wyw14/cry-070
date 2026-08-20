package masking

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

type Snapshot struct {
	ID, SourceID, Digest string
	Rows                 int
	Fields               []string
	ColumnTypes          map[string]string
	Checksums            []string
	SourceVersion        uint64
}

func (s Snapshot) Clone() Snapshot {
	// Deep-copy every nested field so the clone and its source stay independent.
	s.Fields = append([]string(nil), s.Fields...)
	s.Checksums = append([]string(nil), s.Checksums...)
	if s.ColumnTypes != nil {
		columnTypes := make(map[string]string, len(s.ColumnTypes))
		for key, value := range s.ColumnTypes {
			columnTypes[key] = value
		}
		s.ColumnTypes = columnTypes
	}
	return s
}
func (s *Snapshot) Normalize() {
	sort.Strings(s.Fields)
	sort.Strings(s.Checksums)
	if s.ColumnTypes == nil {
		s.ColumnTypes = map[string]string{}
	}
	for key, value := range s.ColumnTypes {
		trimmedKey := strings.TrimSpace(key)
		trimmedValue := strings.TrimSpace(value)
		if trimmedKey == "" {
			delete(s.ColumnTypes, key)
			continue
		}
		if trimmedValue == "" {
			trimmedValue = "text"
		}
		if trimmedKey != key {
			delete(s.ColumnTypes, key)
		}
		s.ColumnTypes[trimmedKey] = trimmedValue
	}
	if s.Rows < 0 {
		s.Rows = 0
	}
}

func (s Snapshot) Validate() error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.SourceID) == "" {
		return ErrInvalid
	}
	if s.Rows < 0 || len(s.Fields) == 0 {
		return ErrInvalid
	}
	for _, field := range s.Fields {
		if strings.TrimSpace(field) == "" {
			return ErrInvalid
		}
	}
	return nil
}

func (s Snapshot) FieldIndex(name string) int {
	name = strings.TrimSpace(name)
	for index, field := range s.Fields {
		if field == name {
			return index
		}
	}
	return -1
}

func (s Snapshot) DigestValue() string {
	hash := sha256.New()
	hash.Write([]byte(s.ID))
	hash.Write([]byte{0})
	hash.Write([]byte(s.SourceID))
	hash.Write([]byte{0})
	for _, field := range s.Fields {
		hash.Write([]byte(field))
		hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil))
}
