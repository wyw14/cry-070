package masking

import "sort"

type ExecutionReport struct {
	BatchID           string
	State             BatchState
	Processed, Failed int
	Differences       []CellDiff
	Warnings          []string
}

func (r ExecutionReport) Success() bool { return r.State == BatchSucceeded && r.Failed == 0 }
func (r *ExecutionReport) AddWarning(value string) {
	if value != "" {
		r.Warnings = append(r.Warnings, value)
		sort.Strings(r.Warnings)
	}
}
func (r ExecutionReport) Summary() map[string]int {
	return map[string]int{"processed": r.Processed, "failed": r.Failed, "differences": len(r.Differences)}
}
