package task

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadTask(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "contenu")
	}))
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "out.txt")
	dt := NewDownloadTask("d", srv.URL, dest)
	if err := dt.Execute(context.Background()); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "contenu" {
		t.Errorf("contenu = %q, attendu %q", string(data), "contenu")
	}
}

func TestDownloadTaskHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	dt := NewDownloadTask("d", srv.URL, filepath.Join(t.TempDir(), "x"))
	if err := dt.Execute(context.Background()); err == nil {
		t.Fatal("attendu une erreur sur un statut HTTP 404")
	}
}
