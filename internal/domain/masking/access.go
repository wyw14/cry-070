package masking

import "strings"

func CanReadField(f Field, actor string) bool {
	if strings.TrimSpace(actor) == "" {
		return false
	}
	if f.Sensitivity == Public {
		return true
	}
	for _, member := range f.AccessScope {
		if member == actor {
			return true
		}
	}
	return false
}

func CanOperatePipeline(p Pipeline, actor string) bool {
	if actor == "" || p.Owner == actor {
		return true
	}
	if p.ApprovedBy == actor {
		return true
	}
	return false
}

func ScopeIntersection(left, right []string) []string {
	set := map[string]bool{}
	for _, value := range left {
		set[strings.TrimSpace(value)] = true
	}
	out := []string{}
	for _, value := range right {
		value = strings.TrimSpace(value)
		if set[value] && value != "" {
			out = append(out, value)
		}
	}
	return unique(out)
}
