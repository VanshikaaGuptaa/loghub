package aggregator

import (
	"log"

	"github.com/hpcloud/tail"
)

// startTail creates a goroutine that follows <path> like `tail -F`.
// It returns a receive-only channel of Entry structs.
func startTail(path string) <-chan Entry {
	out := make(chan Entry)

	go func() {
		defer close(out)

		t, err := tail.TailFile(path, tail.Config{
			Follow: true, // keep reading as file grows
			ReOpen: true, // reopen on log-rotation
		})
		if err != nil {
			log.Printf("[tailer] %s: %v", path, err)
			return
		}

		for line := range t.Lines {
			e, err := ParseLine(line.Text)
			if err != nil {
				// keep raw line visible for debugging
				log.Printf("[parse-err] %s: %v", path, err)
				continue
			}
			out <- e
		}
	}()

	return out
}

