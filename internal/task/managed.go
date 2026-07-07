package task

import "time"

type Managed struct {
	Task
	Timeout time.Duration
	Retries int
}
