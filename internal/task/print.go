package task

import (
	"context"
	"fmt"
	"io"
	"os"
)

type PrintTask struct {
	id      string
	message string
	out     io.Writer
}

func NewPrintTask(id, message string) *PrintTask {
	return &PrintTask{id: id, message: message}
}

func (t *PrintTask) ID() string { return t.id }

func (t *PrintTask) Execute(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return NewTaskError(CodeTimeout, t.id, err)
	}
	out := t.out
	if out == nil {
		out = os.Stderr
	}
	if _, err := fmt.Fprintln(out, t.message); err != nil {
		return NewTaskError(CodeExec, t.id, err)
	}
	return nil
}
