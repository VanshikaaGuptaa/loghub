package aggregator

import (
    "bufio"
    "os"
)

// StreamStdin returns a chan that emits each line from os.Stdin.
func StreamStdin() <-chan string {
    out := make(chan string)
    go func() {
        sc := bufio.NewScanner(os.Stdin)
        for sc.Scan() {
            out <- sc.Text()
        }
        close(out)
    }()
    return out
}
