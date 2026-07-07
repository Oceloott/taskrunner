package metrics

import (
	"fmt"
	"runtime"
	"strings"

	"taskrunner/internal/report"
)

func WriteMetrics(results []report.TaskResult) string {
	var success, failed, timeout int
	for _, r := range results {
		switch r.Status {
		case report.StatusSuccess:
			success++
		case report.StatusFailed:
			failed++
		case report.StatusTimeout:
			timeout++
		}
	}

	var b strings.Builder
	fmt.Fprintln(&b, "# Métriques d'exécution")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "- Goroutines actives à la fin : %d\n", runtime.NumGoroutine())
	fmt.Fprintf(&b, "- Tâches exécutées : %d\n", len(results))
	fmt.Fprintf(&b, "- Tâches réussies : %d\n", success)
	fmt.Fprintf(&b, "- Tâches en échec : %d\n", failed)
	fmt.Fprintf(&b, "- Tâches en timeout : %d\n", timeout)
	return b.String()
}
