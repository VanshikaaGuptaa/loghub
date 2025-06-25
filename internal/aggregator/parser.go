// internal/aggregator/parser.go
package aggregator

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf16"
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
	raw = strings.TrimPrefix(raw, "\uFEFF") // remove UTF-8 BOM
	if strings.HasPrefix(raw, "\xFF\xFE") {            // UTF-16 LE BOM
		b := []byte(raw[2:])                           // skip BOM
		u16 := make([]uint16, len(b)/2)
		for i := 0; i < len(u16); i++ {
			u16[i] = uint16(b[2*i]) | uint16(b[2*i+1])<<8
		}
		raw = string(utf16.Decode(u16))                // ← the missing line
	}
	  // ← remove UTF-8 BOM if present
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
