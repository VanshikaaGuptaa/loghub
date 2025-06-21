package aggregator

import (
	"sync"
)

// StreamMany starts one tailer per file in paths, then returns a single
// channel that emits entries from *all* files, already merged.
func StreamMany(paths []string) <-chan Entry {
	merged := make(chan Entry)
	var wg sync.WaitGroup

	for _, p := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			for e := range startTail(path) { // fan-in
				merged <- e
			}
		}(p)
	}

	// close merged channel *after* all tailers finish (i.e., on CTRL-C)
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}
