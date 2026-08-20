package masking

import "fmt"

type RecoveryDecision struct {
	BatchID  string
	ResumeAt int
	Strategy string
	Safe     bool
}

func PlanRecovery(b Batch, snapshot string) RecoveryDecision {
	d := RecoveryDecision{BatchID: b.ID, ResumeAt: b.Processed, Strategy: "resume-snapshot"}
	if snapshot == "" {
		d.Strategy = "restart"
		d.Safe = false
		return d
	}
	d.Safe = b.State == BatchFailed || b.State == BatchPaused
	return d
}
func ValidateRecovery(d RecoveryDecision, b Batch) error {
	if d.BatchID != b.ID {
		return fmt.Errorf("recovery batch mismatch")
	}
	if d.ResumeAt < 0 || d.ResumeAt > b.Total {
		return fmt.Errorf("recovery offset invalid")
	}
	return nil
}
