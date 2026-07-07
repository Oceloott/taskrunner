package task

import (
	"context"
	"errors"
)

type CalcTask struct {
	id     string
	value  int
	Result int
}

func NewCalcTask(id string, value int) *CalcTask {
	return &CalcTask{id: id, value: value}
}

func (t *CalcTask) ID() string { return t.id }

func (t *CalcTask) Execute(ctx context.Context) error {
	if t.value < 0 {
		return NewTaskError(CodeInvalid, t.id, errors.New("value doit être >= 0"))
	}
	sum := 0
	for i := 1; i <= t.value; i++ {
		if err := ctx.Err(); err != nil {
			return NewTaskError(CodeTimeout, t.id, err)
		}
		sum += i
	}
	t.Result = sum
	return nil
}
