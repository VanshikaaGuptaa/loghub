package aggregator

import (
	

	"github.com/hpcloud/tail"
)

// startTail creates a goroutine that follows <path> like `tail -F`.
// It returns a receive-only channel of Entry structs.
func startTail(path string) (<-chan string, error) {
	t, err := tail.TailFile(path, tail.Config{
		Follow:    true,              // <─ keep the goroutine alive
		ReOpen:    true,    
		Poll:      true,          // <─ reopen after log-rotate
		MustExist: true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // start at end
	})
	if err != nil {
		return nil, err
	}
	out := make(chan string)
	go func() {
		for line := range t.Lines {
			out <- line.Text
		}
		close(out) // shouldn't happen with Follow:true
	}()
	return out, nil
}

