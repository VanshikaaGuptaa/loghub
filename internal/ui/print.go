// internal/ui/print.go
package ui

import (
	"fmt"
	"hash/crc32"
	"strings"
	// "time"
	
    "github.com/mattn/go-colorable"
	"github.com/fatih/color"
	"github.com/VanshikaaGuptaa/loghub/internal/aggregator" // ← use *your* module path
)
func init() {
    // Tell fatih/color to always colour *and* wrap stdout/stderr on Windows
    color.NoColor = false
    color.Output  = colorable.NewColorableStdout()
    color.Error   = colorable.NewColorableStderr()
}

// Print renders one Entry with colours and tidy columns.
func Print(e aggregator.Entry) {
	if e.Level == "" {                // fallback for Raw-only lines
        fmt.Println(e.Raw)
        return
    }
	fmt.Printf("%s %s %-10s %s\n",
		e.Time.Format("15:04:05.000"),      // hh:mm:ss.mmm
		levelColour(e.Level)(strings.ToUpper(e.Level)),
		serviceColour(e.Service)(e.Service),
		e.Msg,
	)
}

////////////////////////////// helpers //////////////////////////////

func levelColour(lvl string) func(a ...any) string {
	switch strings.ToUpper(lvl) {
	case "ERROR", "FATAL":
		return color.New(color.FgRed).SprintFunc()
	case "WARN", "WARNING":
		return color.New(color.FgYellow).SprintFunc()
	case "DEBUG":
		return color.New(color.FgBlue).SprintFunc()
	default: // INFO, TRACE, etc.
		return color.New(color.FgGreen).SprintFunc()
	}
}

func serviceColour(name string) func(a ...any) string {
	// Stable colour picked from hash(name) so "frontend" is always cyan, etc.
	palette := []color.Attribute{
		color.FgCyan, color.FgMagenta, color.FgHiBlue,
		color.FgHiCyan, color.FgHiMagenta, color.FgHiGreen,
	}
	idx := int(crc32.ChecksumIEEE([]byte(name))) % len(palette)
	return color.New(palette[idx]).SprintFunc()
}
