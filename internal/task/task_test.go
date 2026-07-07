package task

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTaskErrorUnwrap(t *testing.T) {
	sentinel := errors.New("boom")
	err := NewTaskError(CodeExec, "t1", sentinel)

	if !errors.Is(err, sentinel) {
		t.Error("errors.Is doit retrouver l'erreur sous-jacente via Unwrap")
	}
	if err.Unwrap() != sentinel {
		t.Error("Unwrap doit retourner l'erreur sous-jacente")
	}
	if !strings.Contains(err.Error(), "t1") {
		t.Errorf("Error() doit contenir l'ID de la tâche, obtenu %q", err.Error())
	}
}

func TestPrintTask(t *testing.T) {
	var buf bytes.Buffer
	pt := &PrintTask{id: "p", message: "coucou", out: &buf}

	if err := pt.Execute(context.Background()); err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if buf.String() != "coucou\n" {
		t.Errorf("sortie = %q, attendu %q", buf.String(), "coucou\n")
	}
}

func TestCalcTask(t *testing.T) {
	ct := NewCalcTask("c", 10)
	if err := ct.Execute(context.Background()); err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if ct.Result != 55 {
		t.Errorf("Result = %d, attendu 55 (somme 1..10)", ct.Result)
	}
}

func TestCalcTaskNegative(t *testing.T) {
	ct := NewCalcTask("c", -1)
	if err := ct.Execute(context.Background()); err == nil {
		t.Fatal("attendu une erreur pour une value négative")
	}
}

func TestFakeTaskSuccess(t *testing.T) {
	ft := NewFakeTask("f", BehaviorSuccess, 5*time.Millisecond)
	if err := ft.Execute(context.Background()); err != nil {
		t.Fatalf("attendu succès, obtenu: %v", err)
	}
}

func TestFakeTaskFail(t *testing.T) {
	ft := NewFakeTask("f", BehaviorFail, 5*time.Millisecond)
	err := ft.Execute(context.Background())
	if err == nil {
		t.Fatal("attendu une erreur, obtenu nil")
	}
	var te *TaskError
	if !errors.As(err, &te) {
		t.Fatalf("attendu *TaskError, obtenu %T", err)
	}
}

func TestFakeTaskTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	ft := NewFakeTask("f", BehaviorTimeout, time.Hour)
	err := ft.Execute(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("attendu DeadlineExceeded, obtenu: %v", err)
	}
}

func TestLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	content := `{"tasks":[
		{"id":"a","type":"print","params":{"message":"hi"},"timeout":"1s","retries":0},
		{"id":"b","type":"calc","params":{"value":5},"timeout":"2s","retries":1},
		{"id":"c","type":"fake","params":{"behavior":"fail","delay":"5ms"},"timeout":"1s","retries":2}
	]}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	tasks, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("len = %d, attendu 3", len(tasks))
	}

	m, ok := tasks[1].(*Managed)
	if !ok {
		t.Fatalf("tasks[1] n'est pas un *Managed (%T)", tasks[1])
	}
	if m.Timeout != 2*time.Second {
		t.Errorf("timeout = %v, attendu 2s", m.Timeout)
	}
	if m.Retries != 1 {
		t.Errorf("retries = %d, attendu 1", m.Retries)
	}
	if _, ok := m.Task.(*CalcTask); !ok {
		t.Errorf("tasks[1] doit envelopper un *CalcTask, obtenu %T", m.Task)
	}
}

func TestLoadUnknownType(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte(`{"tasks":[{"id":"x","type":"zzz","params":{},"timeout":"1s"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("attendu une erreur pour un type inconnu")
	}
	var te *TaskError
	if !errors.As(err, &te) {
		t.Fatalf("attendu *TaskError, obtenu %T", err)
	}
}

func TestLoadInvalidTimeout(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte(`{"tasks":[{"id":"x","type":"calc","params":{"value":1},"timeout":"pas-une-duree"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("attendu une erreur pour un timeout invalide")
	}
}
