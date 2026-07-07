package task

import (
	"context"
	"errors"
	"time"
)

type FakeTaskBehavior int

const (
	BehaviorSuccess FakeTaskBehavior = iota
	BehaviorFail
	BehaviorTimeout
)

type FakeTask struct {
	id       string
	behavior FakeTaskBehavior
	delay    time.Duration
}

func NewFakeTask(id string, behavior FakeTaskBehavior, delay time.Duration) *FakeTask {
	return &FakeTask{id: id, behavior: behavior, delay: delay}
}

func (t *FakeTask) ID() string { return t.id }

func (t *FakeTask) Execute(ctx context.Context) error {
	if t.behavior == BehaviorTimeout {
		<-ctx.Done()
		return NewTaskError(CodeTimeout, t.id, ctx.Err())
	}
	select {
	case <-ctx.Done():
		return NewTaskError(CodeTimeout, t.id, ctx.Err())
	case <-time.After(t.delay):
	}
	if t.behavior == BehaviorFail {
		return NewTaskError(CodeExec, t.id, errors.New("échec simulé"))
	}
	return nil
}
