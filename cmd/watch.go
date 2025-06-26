package cmd

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/VanshikaaGuptaa/loghub/internal/aggregator"
	"github.com/VanshikaaGuptaa/loghub/internal/ui"
	
)

/* ─── flag wiring ──────────────────────────────────────────────────────── */

var (
	watchPath     string
	watchServices []string
	watchExt      string
	useStdin      bool
)

func init() {
	watchCmd.Flags().BoolVar(&useStdin, "stdin", false, "also read from stdin")
	watchCmd.Flags().StringVarP(&watchPath, "path", "p", "./logs", "directory containing log files")
	watchCmd.Flags().StringSliceVarP(&watchServices, "services", "s", nil, "comma-list of service names to include (without extension)")
	watchCmd.Flags().StringVar(&watchExt, "ext", ".log", "log-file extension")
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Stream multiple log files live (plus optional stdin)",
	RunE:  runWatch,
}

/* ─── helper: scan + filter once ───────────────────────────────────────── */

func matchingFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(watchPath, "*"+watchExt))
	if err != nil {
		return nil, err
	}
	if len(watchServices) == 0 {
		return files, nil
	}

	filter := make(map[string]bool, len(watchServices))
	for _, svc := range watchServices {
		filter[svc] = true
	}

	var keep []string
	for _, f := range files {
		svc := strings.TrimSuffix(filepath.Base(f), watchExt)
		if filter[svc] {
			keep = append(keep, f)
		}
	}
	return keep, nil
}

/* ─── main implementation ─────────────────────────────────────────────── */

func runWatch(cmd *cobra.Command, args []string) error {
	// Wait up to 10 s for at least one matching file
	deadline := time.Now().Add(10 * time.Second)

	var files []string
	var err error
	for {
		files, err = matchingFiles()
		if err != nil {
			return err
		}
		if len(files) > 0 || useStdin {
			break // we have something to stream
		}
		if time.Now().After(deadline) {
			log.Printf("no matching log files found after 10 s – exiting")
			return nil
		}
		log.Printf("waiting for log files …")
		time.Sleep(1 * time.Second)
	}

	/* fan-in: tail every file + optional stdin */

	var chans []<-chan string

	for _, f := range files {
		if ch, err := aggregator.StartTail(f); err == nil {
			chans = append(chans, ch)
		} else {
			log.Printf("tail error %s: %v", f, err)
		}
	}
	if useStdin {
		chans = append(chans, aggregator.StreamStdin())
	}

	stream := aggregator.Stream(chans...) // merged channel of Entry

	for e := range stream {
		ui.Print(e)
	}
	return nil
}
