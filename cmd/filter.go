package cmd

import (
	"bufio"
	"io"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var filterFile string

func init() {
	filterCmd.Flags().StringVarP(&filterFile, "file", "f", "", "file to read (defaults to stdin)")
	rootCmd.AddCommand(filterCmd)
}

var filterCmd = &cobra.Command{
	Use:   "filter <pattern>",
	Short: "Filter a log stream or file with a regex / keyword",
	Args:  cobra.ExactArgs(1),
	RunE:  runFilter,
}

func runFilter(cmd *cobra.Command, args []string) error {
	pat := args[0]
	re, err := regexp.Compile("(?i)" + pat) // (?i)=case-insensitive
	if err != nil {
		return err
	}

	var r io.Reader
	if filterFile == "" {
		r = os.Stdin
	} else {
		f, err := os.Open(filterFile)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if re.MatchString(line) {
			os.Stdout.WriteString(line + "\n")
		}
	}
	return sc.Err()
}
