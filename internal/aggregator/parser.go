// internal/aggregator/parser.go
package aggregator

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// 1️⃣  Declare the regular expression at file scope
var lineRE = regexp.MustCompile(`^(\S+)\s+(\w+)\s+(\w+)\s+"([^"]*)"\s*(.*)$`)

var errBadFormat = errors.New("log line does not match expected format")

// 2️⃣  Your Entry struct (or import it here if defined elsewhere)
type Entry struct {
	Time    time.Time
	Level   string
	Service string
	Msg     string
	Fields  map[string]string
	Raw     string
}

// 3️⃣  The ParseLine function everyone else will call
func ParseLine(raw string) (Entry, error) {
	m := lineRE.FindStringSubmatch(raw)
	if m == nil {
		return Entry{Raw: raw}, errBadFormat
	}

	t, err := time.Parse(time.RFC3339Nano, m[1])
	if err != nil {
		return Entry{Raw: raw}, err
	}

	e := Entry{
		Time:    t,
		Level:   strings.ToUpper(m[2]),
		Service: m[3],
		Msg:     m[4],
		Fields:  parseKVs(m[5]),
		Raw:     raw,
	}
	return e, nil
}

// helper: turn key=value key2="x" tail into a map
func parseKVs(tail string) map[string]string {
	if tail == "" {
		return nil
	}
	out := make(map[string]string)
	for _, part := range strings.Fields(tail) {
		if kv := strings.SplitN(part, "=", 2); len(kv) == 2 {
			out[kv[0]] = strings.Trim(kv[1], `"`)
		}
	}
	return out
}
