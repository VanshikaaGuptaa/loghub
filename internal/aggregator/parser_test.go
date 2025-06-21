package aggregator

import "testing"

func TestParseLine(t *testing.T) {
	line := `2025-06-21T14:12:39.456Z  INFO  backend  "Started HTTP server"  port=4000 env=dev`
	e, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Service != "backend"     { t.Errorf("service mismatch") }
	if e.Fields["port"] != "4000" { t.Errorf("port field missing") }
}
