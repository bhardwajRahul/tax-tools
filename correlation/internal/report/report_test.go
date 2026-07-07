package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func sampleReport() Report {
	d := func(s string) time.Time { t, _ := time.Parse("2006-01-02", s); return t }
	return Report{
		Meta: Meta{
			Frequency:  "weekly",
			ReturnKind: "log",
			Start:      d("2024-01-05"),
			End:        d("2024-03-01"),
			Assets: []Asset{
				{Label: "VWRA", Currency: "USD"},
				{Label: "^NSEI", Currency: "INR"},
			},
		},
		Labels: []string{"VWRA", "^NSEI"},
		Correlation: [][]float64{
			{1, 0.5163},
			{0.5163, 1},
		},
		Covariance: [][]float64{{1, 0}, {0, 1}},
		Mean:       []float64{0.0089, 0.0065},
		Stdev:      []float64{0.0080, 0.0085},
		AnnVol:     []float64{0.0579, 0.0614},
		Pairs:      []Pair{{A: "VWRA", B: "^NSEI", R: 0.5163, CI95Lo: -0.2961, CI95Hi: 0.8953}},
		N:          8,
	}
}

// Every row of a Markdown table (header, separator, body) must have the same
// number of runes, otherwise columns won't line up in the raw source.
func TestMarkdownTablesAreAligned(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, sampleReport(), "md"); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(buf.String(), "\n")

	var block []string
	checkBlock := func(rows []string) {
		if len(rows) < 2 {
			return
		}
		want := len([]rune(rows[0]))
		for _, row := range rows {
			if got := len([]rune(row)); got != want {
				t.Errorf("table row widths differ: %d vs %d\n%q", got, want, row)
			}
		}
	}
	for _, ln := range lines {
		if strings.HasPrefix(ln, "|") {
			block = append(block, ln)
			continue
		}
		if len(block) > 0 {
			checkBlock(block)
			block = nil
		}
	}
	checkBlock(block)
}

func TestMarkdownPipeCountsConsistent(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, sampleReport(), "md"); err != nil {
		t.Fatal(err)
	}
	var pipes int
	var inTable bool
	for _, ln := range strings.Split(buf.String(), "\n") {
		if strings.HasPrefix(ln, "|") {
			n := strings.Count(ln, "|")
			if !inTable {
				pipes, inTable = n, true
			} else if n != pipes {
				t.Errorf("inconsistent pipe count: got %d want %d in %q", n, pipes, ln)
			}
		} else {
			inTable = false
		}
	}
}

func TestRenderCSVMatrix(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, sampleReport(), "csv"); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, ",VWRA,^NSEI") {
		t.Errorf("csv header missing labels:\n%s", got)
	}
	if !strings.Contains(got, "VWRA,1,0.5163") {
		t.Errorf("csv row missing expected values:\n%s", got)
	}
}

func TestRenderJSONRoundish(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, sampleReport(), "json"); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	for _, want := range []string{`"frequency": "weekly"`, `"correlation"`, `"observations": 8`} {
		if !strings.Contains(got, want) {
			t.Errorf("json missing %q in:\n%s", want, got)
		}
	}
}

func TestRenderUnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, sampleReport(), "yaml"); err == nil {
		t.Fatal("want error for unknown format")
	}
}
