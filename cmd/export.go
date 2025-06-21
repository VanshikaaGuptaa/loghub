package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/VanshikaaGuptaa/loghub/internal/aggregator"
)

var (
	expPath  string
	expSince time.Duration
	expOut   string
)

func init() {
	exportCmd.Flags().StringVarP(&expPath, "path", "p", "./logs", "directory with log files")
	exportCmd.Flags().DurationVar(&expSince, "since", 10*time.Minute, "window size (e.g. 30m, 2h)")
	exportCmd.Flags().StringVarP(&expOut, "out", "o", "export.json", "output JSON file")
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Write the last <duration> of logs to a JSON file",
	RunE:  runExport,
}

func runExport(cmd *cobra.Command, args []string) error {
	files, err := filepath.Glob(filepath.Join(expPath, "*.log"))
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-expSince)
	var collected []aggregator.Entry

	for _, f := range files {
		collected = append(collected, aggregator.LoadSince(f, cutoff)...)
	}

	sort.Slice(collected, func(i, j int) bool {
		return collected[i].Time.Before(collected[j].Time)
	})

	b, _ := json.MarshalIndent(collected, "", "  ")
	return os.WriteFile(expOut, b, 0644)
}
