package masking

import "sort"

type CellDiff struct{ Field, Before, After string }
type RowDiff struct {
	Key     string
	Changes []CellDiff
}

func CompareRows(before, after []map[string]string, key string) []RowDiff {
	idx := map[string]map[string]string{}
	for _, row := range before {
		idx[row[key]] = row
	}
	out := []RowDiff{}
	for _, row := range after {
		old := idx[row[key]]
		if old == nil {
			continue
		}
		diff := RowDiff{Key: row[key]}
		for name, value := range row {
			if old[name] != value {
				diff.Changes = append(diff.Changes, CellDiff{name, old[name], value})
			}
		}
		if len(diff.Changes) > 0 {
			sort.Slice(diff.Changes, func(i, j int) bool { return diff.Changes[i].Field < diff.Changes[j].Field })
			out = append(out, diff)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}
