package masking

import "sort"

type CellDiff struct{ Field, Before, After string }
type RowDiff struct {
	Key     string
	Status  string // "added", "removed", or "changed"
	Changes []CellDiff
}

func CompareRows(before, after []map[string]string, key string) []RowDiff {
	beforeIdx := map[string]map[string]string{}
	for _, row := range before {
		beforeIdx[row[key]] = row
	}
	afterIdx := map[string]map[string]string{}
	for _, row := range after {
		afterIdx[row[key]] = row
	}

	out := []RowDiff{}

	// Rows present in the target but not in the source were added.
	for _, row := range after {
		k := row[key]
		if beforeIdx[k] != nil {
			continue
		}
		diff := RowDiff{Key: k, Status: "added"}
		for name, value := range row {
			if name == key {
				continue
			}
			diff.Changes = append(diff.Changes, CellDiff{Field: name, Before: "", After: value})
		}
		sort.Slice(diff.Changes, func(i, j int) bool { return diff.Changes[i].Field < diff.Changes[j].Field })
		out = append(out, diff)
	}

	// Rows present in the source but not in the target were removed.
	for _, row := range before {
		k := row[key]
		if afterIdx[k] != nil {
			continue
		}
		diff := RowDiff{Key: k, Status: "removed"}
		for name, value := range row {
			if name == key {
				continue
			}
			diff.Changes = append(diff.Changes, CellDiff{Field: name, Before: value, After: ""})
		}
		sort.Slice(diff.Changes, func(i, j int) bool { return diff.Changes[i].Field < diff.Changes[j].Field })
		out = append(out, diff)
	}

	// Rows present in both whose fields changed.
	for _, row := range after {
		k := row[key]
		old := beforeIdx[k]
		if old == nil {
			continue
		}
		diff := RowDiff{Key: k, Status: "changed"}
		for name, value := range row {
			if name == key {
				continue
			}
			if old[name] != value {
				diff.Changes = append(diff.Changes, CellDiff{Field: name, Before: old[name], After: value})
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
