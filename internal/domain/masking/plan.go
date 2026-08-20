package masking

import "sort"

type ExecutionPlan struct {
	PipelineID string
	Steps      []PlanStep
	Warnings   []string
}
type PlanStep struct {
	FieldID, RuleID string
	Priority        int
}

func BuildPlan(p Pipeline, mappings []Mapping, fields []Field, rules []Rule) ExecutionPlan {
	fieldSet := map[string]bool{}
	for _, f := range fields {
		fieldSet[f.ID] = true
	}
	ruleSet := map[string]Rule{}
	for _, r := range rules {
		ruleSet[r.ID] = r
	}
	out := ExecutionPlan{PipelineID: p.ID}
	for _, m := range mappings {
		if !fieldSet[m.FieldID] {
			out.Warnings = append(out.Warnings, "missing-field:"+m.FieldID)
			continue
		}
		r, ok := ruleSet[m.RuleID]
		if !ok || !r.Active {
			out.Warnings = append(out.Warnings, "inactive-rule:"+m.RuleID)
			continue
		}
		out.Steps = append(out.Steps, PlanStep{m.FieldID, m.RuleID, m.Priority})
	}
	sort.SliceStable(out.Steps, func(i, j int) bool { return out.Steps[i].Priority < out.Steps[j].Priority })
	return out
}
