package orchestrator

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"taskrunner/internal/report"
	"taskrunner/internal/task"
)

func TestValidateWorkers(t *testing.T) {
	cases := []struct {
		in      int
		want    int
		wantErr bool
	}{
		{-5, 3, true},
		{0, 3, true},
		{1, 1, false},
		{50, 50, false},
		{100, 100, false},
		{101, 3, true},
	}
	for _, c := range cases {
		got, err := ValidateWorkers(c.in)
		if got != c.want {
			t.Errorf("ValidateWorkers(%d) = %d, attendu %d", c.in, got, c.want)
		}
		if (err != nil) != c.wantErr {
			t.Errorf("ValidateWorkers(%d) err = %v, wantErr = %v", c.in, err, c.wantErr)
		}
	}
}

func TestOptions(t *testing.T) {
	cfg := OrchestratorConfig{}
	WithWorkers(7)(&cfg)
	WithVerbose(true)(&cfg)

	if cfg.Workers != 7 {
		t.Errorf("Workers = %d, attendu 7", cfg.Workers)
	}
	if !cfg.Verbose {
		t.Error("Verbose devrait être true")
	}
}

func TestOrchestrateStatuses(t *testing.T) {
	tasks := []task.Task{
		&task.Managed{Task: task.NewFakeTask("ok", task.BehaviorSuccess, 0), Timeout: time.Second, Retries: 0},
		&task.Managed{Task: task.NewFakeTask("fail", task.BehaviorFail, 0), Timeout: time.Second, Retries: 1},
		&task.Managed{Task: task.NewFakeTask("to", task.BehaviorTimeout, 0), Timeout: 30 * time.Millisecond, Retries: 1},
	}

	rep, err := Orchestrate(context.Background(), tasks, 2)
	if err != nil {
		t.Fatalf("Orchestrate: %v", err)
	}

	byID := map[string]report.TaskResult{}
	for _, r := range rep.Results {
		byID[r.ID] = r
	}

	if r := byID["ok"]; r.Status != report.StatusSuccess || r.Attempts != 1 {
		t.Errorf("ok: status=%s attempts=%d, attendu success/1", r.Status, r.Attempts)
	}
	if r := byID["fail"]; r.Status != report.StatusFailed || r.Attempts != 2 {
		t.Errorf("fail: status=%s attempts=%d, attendu failed/2", r.Status, r.Attempts)
	}
	if r := byID["to"]; r.Status != report.StatusTimeout || r.Attempts != 2 {
		t.Errorf("to: status=%s attempts=%d, attendu timeout/2", r.Status, r.Attempts)
	}
}

type countingTask struct {
	id      string
	current *int32
	max     *int32
}

func (c *countingTask) ID() string { return c.id }

func (c *countingTask) Execute(ctx context.Context) error {
	n := atomic.AddInt32(c.current, 1)
	for {
		m := atomic.LoadInt32(c.max)
		if n <= m || atomic.CompareAndSwapInt32(c.max, m, n) {
			break
		}
	}
	time.Sleep(20 * time.Millisecond)
	atomic.AddInt32(c.current, -1)
	return nil
}

func TestOrchestrateWorkerLimit(t *testing.T) {
	var current, max int32
	var tasks []task.Task
	for i := 0; i < 12; i++ {
		tasks = append(tasks, &task.Managed{
			Task:    &countingTask{id: fmt.Sprint(i), current: &current, max: &max},
			Timeout: time.Second,
		})
	}

	if _, err := Orchestrate(context.Background(), tasks, 3); err != nil {
		t.Fatalf("Orchestrate: %v", err)
	}
	if max > 3 {
		t.Errorf("concurrence max observée = %d, ne doit pas dépasser 3 workers", max)
	}
}
