package cli

import (
	"fmt"
	"github.com/20xygen/git-blame/internal/format"
	"github.com/20xygen/git-blame/internal/statistics"
	"github.com/20xygen/git-blame/internal/utils"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path/filepath"
)

var (
	rootCmd = &cobra.Command{
		Use:   "blame",
		Short: "Analyze git blame statistics",
		Long:  "Blame is a CLI tool to analyze git blame statistics with different formats and features.",
		Args:  cobra.NoArgs,
		Run:   command,
	}
)

func command(cmd *cobra.Command, _ []string) {
	ps, err := statistics.GetParams(*cmd)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		slog.Error(err.Error())
		os.Exit(utils.CodeParametersParsing)
		return
	}

	logger := utils.SetupLogger()
	slog.SetDefault(logger)

	ps.Path, err = filepath.Abs(ps.Path)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		slog.Error(err.Error())
		os.Exit(utils.CodeAbsolutePath)
		return
	}

	info, err := utils.GetLangInfo()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		slog.Error(err.Error())
		os.Exit(utils.CodeLanguageInfo)
		return
	}

	st, err := statistics.CollectStat(ps, info)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		slog.Error(err.Error())
		os.Exit(3)
		return
	}

	output, err := format.AutoFormat(st, ps.OrderBy, ps.Format)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		slog.Error(err.Error())
		os.Exit(utils.CodeFormat)
	}

	fmt.Print(output)

	slog.Info("Done successfully")
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringP("repository", "r", ".", "Git repository path")
	rootCmd.Flags().StringP("revision", "R", "HEAD", "Git revision")
	rootCmd.Flags().StringSliceP("order-by", "o", []string{"lines", "commits", "files"}, "Sort key as comma-separated list of 'lines', 'commits', 'names' or 'files'")
	rootCmd.Flags().BoolP("use-committer", "C", false, "Use committer instead of author")
	rootCmd.Flags().StringSliceP("extensions", "e", nil, "File extensions filter (comma-separated)")
	rootCmd.Flags().StringSliceP("languages", "l", nil, "Languages filter (comma-separated)")
	rootCmd.Flags().StringSliceP("exclude", "x", nil, "Exclude glob patterns")
	rootCmd.Flags().StringSliceP("restrict-to", "t", nil, "Restrict-to glob patterns")
	rootCmd.Flags().StringP("format", "f", "tabular", "Output format (one of 'pretty', 'tabular', 'json', 'json-lines', 'csv')'")
}
