package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// MCMeta describes a Monte-Carlo run.
type MCMeta struct {
	Symbol   string   `json:"symbol"`
	Strategy string   `json:"strategy"`
	Start    string   `json:"start"`
	End      string   `json:"end"`
	Bars     int      `json:"bars"`
	Trials   int      `json:"trials"`
	Seed     int64    `json:"seed"`
	Notes    []string `json:"notes,omitempty"`
}

// MCMetric is one metric's observed value and its bootstrap distribution.
type MCMetric struct {
	Name    string  `json:"name"`
	Percent bool    `json:"-"` // render as % (returns/drawdown) vs plain (Sharpe)
	Actual  float64 `json:"actual"`
	P5      float64 `json:"p5"`
	P25     float64 `json:"p25"`
	P50     float64 `json:"p50"`
	P75     float64 `json:"p75"`
	P95     float64 `json:"p95"`
}

// MonteCarlo is the full bootstrap report.
type MonteCarlo struct {
	Meta    MCMeta     `json:"meta"`
	Metrics []MCMetric `json:"metrics"`
}

// RenderMonteCarlo writes mc to w as "md", "csv" or "json".
func RenderMonteCarlo(w io.Writer, mc MonteCarlo, format string) error {
	switch format {
	case "md", "markdown", "":
		return renderMCMarkdown(w, mc)
	case "csv":
		return renderMCCSV(w, mc)
	case "json":
		return renderMCJSON(w, mc)
	default:
		return fmt.Errorf("report: unknown format %q (want md|csv|json)", format)
	}
}

func renderMCMarkdown(w io.Writer, mc MonteCarlo) error {
	var b strings.Builder
	m := mc.Meta
	fmt.Fprintf(&b, "# Monte-Carlo — %s (%s)\n\n", m.Symbol, m.Strategy)
	fmt.Fprintf(&b, "- Period: %s → %s (%d bars)\n", m.Start, m.End, m.Bars)
	fmt.Fprintf(&b, "- Bootstrap trials: %d (seed %d)\n\n", m.Trials, m.Seed)

	header := []string{"Metric", "Actual", "P5", "P25", "Median", "P75", "P95"}
	aligns := []mdAlign{alignLeft, alignRight, alignRight, alignRight, alignRight, alignRight, alignRight}
	rows := make([][]string, 0, len(mc.Metrics))
	for _, mt := range mc.Metrics {
		f := num
		if mt.Percent {
			f = pct
		}
		rows = append(rows, []string{mt.Name, f(mt.Actual), f(mt.P5), f(mt.P25), f(mt.P50), f(mt.P75), f(mt.P95)})
	}
	mdTable(&b, header, rows, aligns)

	b.WriteByte('\n')
	for _, n := range m.Notes {
		fmt.Fprintf(&b, "> %s\n", n)
	}
	if len(m.Notes) > 0 {
		b.WriteByte('\n')
	}
	fmt.Fprintf(&b, "_%s_\n", disclaimer)

	_, err := io.WriteString(w, b.String())
	return err
}

func renderMCCSV(w io.Writer, mc MonteCarlo) error {
	cw := csv.NewWriter(w)
	rows := [][]string{{"metric", "actual", "p5", "p25", "p50", "p75", "p95"}}
	for _, mt := range mc.Metrics {
		rows = append(rows, []string{
			mt.Name, ff(mt.Actual), ff(mt.P5), ff(mt.P25), ff(mt.P50), ff(mt.P75), ff(mt.P95),
		})
	}
	if err := cw.WriteAll(rows); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}

func renderMCJSON(w io.Writer, mc MonteCarlo) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(mc)
}
