package memory

import "sort"

// IdempotencyStats is a read-only summary used by batch dashboards.
type IdempotencyStats struct {
	TotalBatches  int
	UniqueKeys    int
	DuplicateKeys []string
}

func BuildIdempotencyStats(keys map[string]string) IdempotencyStats {
	seen := map[string]bool{}
	duplicates := []string{}
	for key, id := range keys {
		if key == "" || id == "" {
			continue
		}
		if seen[key] {
			duplicates = append(duplicates, key)
		}
		seen[key] = true
	}
	sort.Strings(duplicates)
	return IdempotencyStats{TotalBatches: len(keys), UniqueKeys: len(seen), DuplicateKeys: duplicates}
}

func (s IdempotencyStats) Healthy() bool {
	return len(s.DuplicateKeys) == 0 && s.UniqueKeys == s.TotalBatches
}

func (s IdempotencyStats) Labels() []string {
	labels := []string{"batches", "keys"}
	if len(s.DuplicateKeys) > 0 {
		labels = append(labels, "duplicates")
	}
	return labels
}
