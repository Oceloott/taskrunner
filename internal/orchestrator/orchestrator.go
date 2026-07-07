package orchestrator

import (
	"context"
	"errors"
	"sync"
	"time"

	"taskrunner/internal/report"
	"taskrunner/internal/task"
)

const defaultTimeout = 30 * time.Second

func Orchestrate(ctx context.Context, tasks []task.Task, workers int, opts ...Option) (report.Report, error) {
	cfg := OrchestratorConfig{Workers: workers}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.Workers < 1 {
		cfg.Workers = 1
	}

	results := make([]report.TaskResult, len(tasks))
	sem := make(chan struct{}, cfg.Workers)
	var wg sync.WaitGroup

	for i, t := range tasks {
		wg.Add(1)
		go func(i int, t task.Task) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results[i] = report.TaskResult{
					ID:       t.ID(),
					Status:   report.StatusTimeout,
					Duration: time.Duration(0).String(),
					Attempts: 0,
				}
				return
			}
			defer func() { <-sem }()
			results[i] = runTask(ctx, t)
		}(i, t)
	}

	wg.Wait()
	return report.Report{Results: results}, nil
}

func runTask(ctx context.Context, t task.Task) report.TaskResult {
	timeout, retries := policy(t)
	start := time.Now()
	status := report.StatusFailed
	attempts := 0

	for attempts <= retries {
		if ctx.Err() != nil {
			status = report.StatusTimeout
			break
		}
		attempts++

		tctx, cancel := context.WithTimeout(ctx, timeout)
		err := t.Execute(tctx)
		cancel()

		if err == nil {
			status = report.StatusSuccess
			break
		}
		if errors.Is(err, context.DeadlineExceeded) {
			status = report.StatusTimeout
		} else {
			status = report.StatusFailed
		}
	}

	return report.TaskResult{
		ID:       t.ID(),
		Status:   status,
		Duration: time.Since(start).String(),
		Attempts: attempts,
	}
}

func policy(t task.Task) (time.Duration, int) {
	if m, ok := t.(*task.Managed); ok {
		timeout := m.Timeout
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		return timeout, m.Retries
	}
	return defaultTimeout, 0
}
