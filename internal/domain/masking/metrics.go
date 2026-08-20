package masking

import "sort"

type Metrics struct {
	Sources, Fields, Rules, Pipelines, Batches int
	BySensitivity                              map[Sensitivity]int
	ByState                                    map[BatchState]int
}

func CollectMetrics(sources []Source, fields []Field, rules []Rule, pipelines []Pipeline, batches []Batch) Metrics {
	m := Metrics{Sources: len(sources), Fields: len(fields), Rules: len(rules), Pipelines: len(pipelines), Batches: len(batches), BySensitivity: map[Sensitivity]int{}, ByState: map[BatchState]int{}}
	for _, f := range fields {
		m.BySensitivity[f.Sensitivity]++
	}
	for _, b := range batches {
		m.ByState[b.State]++
	}
	return m
}
func (m Metrics) Sensitivities() []Sensitivity {
	out := make([]Sensitivity, 0, len(m.BySensitivity))
	for k := range m.BySensitivity {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
func (m Metrics) States() []BatchState {
	out := make([]BatchState, 0, len(m.ByState))
	for k := range m.ByState {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
