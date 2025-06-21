// internal/aggregator/load.go
package aggregator

import (
	"bufio"
	"os"
	"time"
)

// LoadSince reads <path> once and returns every Entry whose timestamp
// is ≥ cutoff.  It uses ParseLine from parser.go.
func LoadSince(path string, cutoff time.Time) []Entry {
	f, err := os.Open(path)
	if err != nil {
		// Return nil on error so callers can still continue.
		return nil
	}
	defer f.Close()

	var out []Entry
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		e, err := ParseLine(sc.Text())
		if err == nil && !e.Time.Before(cutoff) {
			out = append(out, e)
		}
	}
	return out
}
