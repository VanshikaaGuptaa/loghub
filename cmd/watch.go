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

/* -------- flag wiring -------- */

var (
	watchPath     string
	watchServices []string
	watchExt      string
)

func init() {
	watchCmd.Flags().StringVarP(&watchPath, "path", "p", "./logs", "directory containing log files")
	watchCmd.Flags().StringSliceVarP(&watchServices, "services", "s", nil, "comma-list of service names to include (without extension)")
	watchCmd.Flags().StringVar(&watchExt, "ext", ".log", "log-file extension")
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Stream multiple log files live",
	RunE:  runWatch,
}

/* -------- helper: scan + filter once -------- */

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

/* -------- main implementation -------- */

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
		if len(files) > 0 {
			break // success: we have files to tail
		}
		if time.Now().After(deadline) {
			log.Printf("no matching log files found after 10 s – exiting")
			return nil
		}
		log.Printf("waiting for log files …")
		time.Sleep(1 * time.Second)
	}

	// Fan-in tailers and stream to UI
	stream := aggregator.StreamMany(files)
	for e := range stream {
		ui.Print(e)
	}
	return nil
}
