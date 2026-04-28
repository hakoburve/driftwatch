package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Format controls the output format for drift results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Formatter writes drift results to a writer in a specified format.
type Formatter struct {
	Writer io.Writer
	Format Format
}

// New creates a new Formatter with the given writer and format.
func New(w io.Writer, f Format) *Formatter {
	return &Formatter{Writer: w, Format: f}
}

// Write outputs the drift results according to the configured format.
func (f *Formatter) Write(results []drift.Result) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(results)
	default:
		return f.writeText(results)
	}
}

func (f *Formatter) writeText(results []drift.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(f.Writer, "✓ No drift detected.")
		return err
	}
	for _, r := range results {
		status := "DRIFTED"
		if r.Missing {
			status = "MISSING"
		}
		fmt.Fprintf(f.Writer, "[%s] %s\n", status, r.Service)
		for _, d := range r.Diffs {
			fmt.Fprintf(f.Writer, "  key=%s expected=%q actual=%q\n", d.Key, d.Expected, d.Actual)
		}
		if len(r.ExtraKeys) > 0 {
			fmt.Fprintf(f.Writer, "  extra keys: %s\n", strings.Join(r.ExtraKeys, ", "))
		}
	}
	return nil
}

func (f *Formatter) writeJSON(results []drift.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(f.Writer, `{"status":"ok","drift":[]}`)
		return err
	}
	fmt.Fprintln(f.Writer, `{"status":"drifted","drift":["`)
	for i, r := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(f.Writer, `  {"service":%q,"missing":%v,"diffs":%d,"extra_keys":%d}%s\n`,
			r.Service, r.Missing, len(r.Diffs), len(r.ExtraKeys), comma)
	}
	_, err := fmt.Fprintln(f.Writer, "]}")  
	return err
}
