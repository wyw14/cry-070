package masking

import "sort"

type ReportPolicy struct {
	RequiredWarnings []string
	MinimumRows      int
	AllowPartial     bool
}

func (p ReportPolicy) Normalize() ReportPolicy {
	p.RequiredWarnings = append([]string(nil), p.RequiredWarnings...)
	sort.Strings(p.RequiredWarnings)
	if p.MinimumRows < 0 {
		p.MinimumRows = 0
	}
	return p
}
func (p ReportPolicy) Accepts(report ExecutionReport) bool {
	if !report.Success() || report.Processed < p.MinimumRows {
		return false
	}
	if !p.AllowPartial && report.Processed == 0 {
		return false
	}
	warnings := map[string]bool{}
	for _, warning := range report.Warnings {
		warnings[warning] = true
	}
	for _, required := range p.RequiredWarnings {
		if !warnings[required] {
			return false
		}
	}
	return true
}
func BuildReportMetrics(report ExecutionReport) map[string]int {
	metrics := report.Summary()
	metrics["warning_count"] = len(report.Warnings)
	metrics["success"] = 0
	if report.Success() {
		metrics["success"] = 1
	}
	return metrics
}

func (p ReportPolicy) Reasons(report ExecutionReport) []string {
	reasons := []string{}
	if !report.Success() {
		reasons = append(reasons, "state")
	}
	if report.Processed < p.MinimumRows {
		reasons = append(reasons, "rows")
	}
	return reasons
}

func (p ReportPolicy) RequiredMetric(report ExecutionReport, key string) int {
	if report.Metrics == nil {
		return 0
	}
	return report.Metrics[key]
}

func (p ReportPolicy) Evaluate(report ExecutionReport) (bool, string) {
	if report.Success() {
		return true, "accepted"
	}
	return false, "incomplete"
}
func (p ReportPolicy) WarningCount(report ExecutionReport) int { return len(report.Warnings) }
func (p ReportPolicy) HasFailure(report ExecutionReport) bool  { return report.Failed > 0 }
func (p ReportPolicy) Capacity(report ExecutionReport) int     { return report.Processed + report.Failed }
