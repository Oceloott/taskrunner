package metrics

import (
	"strings"
	"testing"

	"taskrunner/internal/report"
)

func TestWriteMetrics(t *testing.T) {
	results := []report.TaskResult{
		{ID: "a", Status: report.StatusSuccess},
		{ID: "b", Status: report.StatusFailed},
		{ID: "c", Status: report.StatusTimeout},
		{ID: "d", Status: report.StatusSuccess},
	}

	out := WriteMetrics(results)

	wants := []string{
		"Tâches exécutées : 4",
		"Tâches réussies : 2",
		"Tâches en échec : 1",
		"Tâches en timeout : 1",
		"Goroutines actives",
	}
	for _, w := range wants {
		if !strings.Contains(out, w) {
			t.Errorf("METRICS doit contenir %q\n--- contenu obtenu ---\n%s", w, out)
		}
	}
}
