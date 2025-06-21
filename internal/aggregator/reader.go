package aggregator

import (
	"bufio"
	"fmt"
	"os"
)

// StreamFile opens path, scans it line-by-line, and sends Entry structs
// into the returned channel.  Closes the channel when EOF is reached.
func StreamFile(path string) (<-chan Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	out := make(chan Entry)

	go func() {
		defer f.Close()
		defer close(out)

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			raw := scanner.Text()

			entry, err := ParseLine(raw) // <- the parser you already wrote
			if err != nil {
				// Could log the error or just forward the raw text:
				fmt.Println(raw)
				continue
			}
			out <- entry
		}
	}()

	return out, nil
}
