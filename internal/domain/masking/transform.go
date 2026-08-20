package masking

import "strings"

type TransformResult struct {
	Value   string
	Changed bool
	RuleID  string
	Error   string
}

func ApplyRule(rule Rule, value string) TransformResult {
	out, err := rule.Validate(value)
	if err != nil {
		return TransformResult{Value: value, RuleID: rule.ID, Error: err.Error()}
	}
	return TransformResult{Value: out, Changed: out != value, RuleID: rule.ID}
}
func ApplyPlan(plan ExecutionPlan, rules map[string]Rule, values map[string]string) map[string]TransformResult {
	out := map[string]TransformResult{}
	for _, step := range plan.Steps {
		r, ok := rules[step.RuleID]
		if !ok {
			out[step.FieldID] = TransformResult{Value: values[step.FieldID], Error: "rule missing"}
			continue
		}
		out[step.FieldID] = ApplyRule(r, strings.TrimSpace(values[step.FieldID]))
	}
	return out
}
