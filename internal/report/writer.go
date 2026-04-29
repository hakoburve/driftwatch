package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// WriteText writes a human-readable report to w.
func WriteText(w io.Writer, r Report) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Generated:\t%s\n", r.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(tw, "Status:\t%s\n", r.Status)
	fmt.Fprintf(tw, "Total:\t%d\tDrifted:\t%d\tMissing:\t%d\tIn-Sync:\t%d\n",
		r.Summary.Total, r.Summary.Drifted, r.Summary.Missing, r.Summary.InSync)
	if err := tw.Flush(); err != nil {
		return err
	}
	for _, res := range r.Results {
		if err := writeTextResult(w, res); err != nil {
			return err
		}
	}
	return nil
}

func writeTextResult(w io.Writer, res drift.Result) error {
	if res.Missing {
		_, err := fmt.Fprintf(w, "\n[MISSING] %s\n", res.Service)
		return err
	}
	if len(res.Diffs) == 0 {
		_, err := fmt.Fprintf(w, "\n[OK]      %s\n", res.Service)
		return err
	}
	fmt.Fprintf(w, "\n[DRIFT]   %s\n", res.Service)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "  KEY\tDECLARED\tLIVE")
	for _, d := range res.Diffs {
		fmt.Fprintf(tw, "  %s\t%s\t%s\n", d.Key, d.Declared, d.Live)
	}
	return tw.Flush()
}

// WriteJSON writes a machine-readable JSON report to w.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
