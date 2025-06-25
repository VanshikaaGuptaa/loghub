package aggregator

import (
	"log"
	"sync"
)

// StreamMany starts one tail-goroutine per file, parses each line,
// and fans everything into a single channel of Entry.
//
// • The output channel is never closed, so a brief tail restart on Windows
//   can't terminate the whole program.
//
// • If ParseLine() returns an error the raw line is forwarded anyway, so
//   unstructured messages still appear in the UI.
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

			for line := range lines {
				e, err := ParseLine(line)
				if err != nil {
					// Unknown format – keep raw text so user still sees it
					out <- Entry{Raw: line}
					continue
				}
				out <- e
			}
		}(p)
	}

	// Prevent wg leak, but deliberately never close(out)
	go func() { wg.Wait() }()

	return out
}
