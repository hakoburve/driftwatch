package drift

// Summary holds aggregate counts across all checked services.
type Summary struct {
	Total   int `json:"total"`
	InSync  int `json:"in_sync"`
	Drifted int `json:"drifted"`
	Missing int `json:"missing"`
}

// Summarize computes a Summary from a slice of Results.
func Summarize(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch {
		case r.Missing:
			s.Missing++
		case len(r.Diffs) > 0:
			s.Drifted++
		default:
			s.InSync++
		}
	}
	return s
}
