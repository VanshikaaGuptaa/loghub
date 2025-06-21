package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/VanshikaaGuptaa/loghub/internal/aggregator"
	"github.com/VanshikaaGuptaa/loghub/internal/ui"
)

// --- flags added in init() ---
var (
	watchPath     string
	watchServices []string
	watchExt      string
)

func init() {
	watchCmd.Flags().StringVarP(&watchPath, "path", "p", "./logs", "directory containing log files")
	watchCmd.Flags().StringSliceVarP(&watchServices, "services", "s", nil, "comma-list of service names to include")
	watchCmd.Flags().StringVar(&watchExt, "ext", ".log", "log-file extension")
	rootCmd.AddCommand(watchCmd)
}

// watchCmd definition (kept from cobra-cli)
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Stream multiple log files live",
	RunE:  runWatch,
}

// ------------ actual implementation -------------
func runWatch(cmd *cobra.Command, args []string) error {

	files, err := filepath.Glob(filepath.Join(watchPath, "*"+watchExt))
	if err != nil {
		return err
	}
	if len(watchServices) > 0 {
		// keep only svc names the user asked for
		filter := map[string]bool{}
		for _, svc := range watchServices {
			filter[svc] = true
		}
		var keep []string
		for _, f := range files {
			base := filepath.Base(f)           // e.g. backend.log
			svc := strings.TrimSuffix(base, watchExt)
			if filter[svc] {
				keep = append(keep, f)
			}
		}
		files = keep
	}
	if len(files) == 0 {
		log.Println("no matching log files found")
		return nil
	}

	stream := aggregator.StreamMany(files)
	for e := range stream {
		ui.Print(e)
	}
	return nil
}
