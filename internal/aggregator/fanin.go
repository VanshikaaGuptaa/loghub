package aggregator

import (
	"log"
	"sync"
)

// StreamMany starts one tail-goroutine per file, parses each line,
// and fans everything into a single channel of Entry.
// The output channel is intentionally *never* closed so that a brief
// tail restart on Windows can't terminate the whole program.
func StreamMany(paths []string) <-chan Entry {
	out := make(chan Entry)
	var wg sync.WaitGroup

	for _, p := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			lines, err := startTail(path)
			if err != nil {
				log.Printf("tail error %s: %v", path, err)
				return
			}

			for line := range lines {           // this channel may re-appear
				e, err := ParseLine(line)
				if err == nil {
					out <- e
				}
			}
		}(p)
	}

	// keep a background goroutine so we don't leak wg,
	// but DON'T close(out) when wg reaches zero
	go func() { wg.Wait() }()

	return out
}
