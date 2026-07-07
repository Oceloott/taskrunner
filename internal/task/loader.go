package task

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const defaultTimeout = 30 * time.Second

type rawTask struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Params  json.RawMessage `json:"params"`
	Timeout string          `json:"timeout"`
	Retries int             `json:"retries"`
}

type rawFile struct {
	Tasks []rawTask `json:"tasks"`
}

func Load(path string) ([]Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lecture du fichier %q: %w", path, err)
	}
	var rf rawFile
	if err := json.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	tasks := make([]Task, 0, len(rf.Tasks))
	for _, rt := range rf.Tasks {
		base, err := build(rt)
		if err != nil {
			return nil, err
		}
		timeout := defaultTimeout
		if rt.Timeout != "" {
			timeout, err = time.ParseDuration(rt.Timeout)
			if err != nil {
				return nil, NewTaskError(CodeLoad, rt.ID, fmt.Errorf("timeout invalide %q: %w", rt.Timeout, err))
			}
		}
		tasks = append(tasks, &Managed{Task: base, Timeout: timeout, Retries: rt.Retries})
	}
	return tasks, nil
}

func build(rt rawTask) (Task, error) {
	switch rt.Type {
	case "print":
		var p struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(rt.Params, &p); err != nil {
			return nil, NewTaskError(CodeLoad, rt.ID, err)
		}
		return NewPrintTask(rt.ID, p.Message), nil

	case "calc":
		var p struct {
			Value int `json:"value"`
		}
		if err := json.Unmarshal(rt.Params, &p); err != nil {
			return nil, NewTaskError(CodeLoad, rt.ID, err)
		}
		return NewCalcTask(rt.ID, p.Value), nil

	case "download":
		var p struct {
			URL  string `json:"url"`
			Dest string `json:"dest"`
		}
		if err := json.Unmarshal(rt.Params, &p); err != nil {
			return nil, NewTaskError(CodeLoad, rt.ID, err)
		}
		return NewDownloadTask(rt.ID, p.URL, p.Dest), nil

	case "fake":
		var p struct {
			Behavior string `json:"behavior"`
			Delay    string `json:"delay"`
		}
		if err := json.Unmarshal(rt.Params, &p); err != nil {
			return nil, NewTaskError(CodeLoad, rt.ID, err)
		}
		behavior, err := parseBehavior(p.Behavior)
		if err != nil {
			return nil, NewTaskError(CodeLoad, rt.ID, err)
		}
		var delay time.Duration
		if p.Delay != "" {
			delay, err = time.ParseDuration(p.Delay)
			if err != nil {
				return nil, NewTaskError(CodeLoad, rt.ID, fmt.Errorf("delay invalide %q: %w", p.Delay, err))
			}
		}
		return NewFakeTask(rt.ID, behavior, delay), nil

	default:
		return nil, NewTaskError(CodeLoad, rt.ID, fmt.Errorf("type de tâche inconnu %q", rt.Type))
	}
}

func parseBehavior(s string) (FakeTaskBehavior, error) {
	switch s {
	case "success", "":
		return BehaviorSuccess, nil
	case "fail":
		return BehaviorFail, nil
	case "timeout":
		return BehaviorTimeout, nil
	default:
		return BehaviorSuccess, fmt.Errorf("comportement fake inconnu %q", s)
	}
}
