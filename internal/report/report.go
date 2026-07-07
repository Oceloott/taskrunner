package report

import (
	"encoding/json"
	"io"
)

const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
	StatusTimeout = "timeout"
)

type TaskResult struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Duration string `json:"duration"`
	Attempts int    `json:"attempts"`
}

type Report struct {
	Results []TaskResult `json:"results"`
}

func (r Report) WriteTo(w io.Writer) (int64, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return 0, err
	}
	data = append(data, '\n')
	n, err := w.Write(data)
	return int64(n), err
}
