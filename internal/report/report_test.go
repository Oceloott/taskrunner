package report

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestReportWriteTo(t *testing.T) {
	r := Report{Results: []TaskResult{
		{ID: "t1", Status: StatusSuccess, Duration: "12ms", Attempts: 1},
		{ID: "t2", Status: StatusTimeout, Duration: "3.001s", Attempts: 3},
	}}

	var buf bytes.Buffer
	n, err := r.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo: %v", err)
	}
	if n != int64(buf.Len()) {
		t.Errorf("n = %d, mais buffer = %d octets", n, buf.Len())
	}

	var got Report
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("la sortie n'est pas du JSON valide: %v", err)
	}
	if len(got.Results) != 2 {
		t.Fatalf("len = %d, attendu 2", len(got.Results))
	}
	if got.Results[1].ID != "t2" || got.Results[1].Status != StatusTimeout {
		t.Errorf("round-trip incorrect: %+v", got.Results[1])
	}
}
