// Package task définit l'interface Task, ses implémentations et l'erreur custom TaskError.
package task

import (
	"context"
	"fmt"
)

type Task interface {
	ID() string
	Execute(ctx context.Context) error
}

const (
	CodeUnknown = iota
	CodeExec
	CodeTimeout
	CodeLoad
	CodeInvalid
)

type TaskError struct {
	Code   int
	TaskID string
	Err    error
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("task %q (code %d): %v", e.TaskID, e.Code, e.Err)
}

func (e *TaskError) Unwrap() error {
	return e.Err
}

func NewTaskError(code int, taskID string, err error) *TaskError {
	return &TaskError{Code: code, TaskID: taskID, Err: err}
}
