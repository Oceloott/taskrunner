package task

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DownloadTask struct {
	id   string
	url  string
	dest string
}

func NewDownloadTask(id, url, dest string) *DownloadTask {
	return &DownloadTask{id: id, url: url, dest: dest}
}

func (t *DownloadTask) ID() string { return t.id }

func (t *DownloadTask) Execute(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.url, nil)
	if err != nil {
		return NewTaskError(CodeInvalid, t.id, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return NewTaskError(CodeExec, t.id, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return NewTaskError(CodeExec, t.id, fmt.Errorf("statut HTTP %d", resp.StatusCode))
	}
	f, err := os.Create(t.dest)
	if err != nil {
		return NewTaskError(CodeExec, t.id, err)
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return NewTaskError(CodeExec, t.id, err)
	}
	return nil
}
