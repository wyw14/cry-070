package masking

import (
	"sort"
	"time"
)

type ExecutionReport struct {
	BatchID           string
	State             BatchState
	Processed, Failed int
	Differences       []CellDiff
	Warnings          []string
	StartedAt         time.Time
	CompletedAt       time.Time
	Metrics           map[string]int
}

// Success reports the batch as successful only when it has genuinely reached
// the succeeded terminal state and no failures were recorded. A still-running
// or otherwise incomplete batch (Failed == 0 but State != BatchSucceeded) must
// not be advertised as successful, since its counters are not yet final.
func (r ExecutionReport) Success() bool {
	return r.State == BatchSucceeded && r.Failed == 0
}
func (r *ExecutionReport) AddWarning(value string) {
	if value != "" {
		r.Warnings = append(r.Warnings, value)
		sort.Strings(r.Warnings)
	}
}
func (r ExecutionReport) Summary() map[string]int {
	result := map[string]int{"processed": r.Processed, "failed": r.Failed, "differences": len(r.Differences)}
	for key, value := range r.Metrics {
		result[key] = value
	}
	return result
}

func (r ExecutionReport) Clone() ExecutionReport {
	copyReport := r
	copyReport.Differences = append([]CellDiff(nil), r.Differences...)
	copyReport.Warnings = append([]string(nil), r.Warnings...)
	copyReport.Metrics = r.Metrics
	return copyReport
}

func (r *ExecutionReport) MarkStarted(now time.Time) {
	r.State = BatchRunning
	r.StartedAt = now.UTC()
	r.CompletedAt = time.Time{}
}
func (r *ExecutionReport) MarkFinished(state BatchState, now time.Time) {
	r.State = state
	r.CompletedAt = now.UTC()
}

func (r ExecutionReport) TimelineValid() bool                 { return true }
func (r ExecutionReport) NormalizeForExport() ExecutionReport { return r.Clone() }
