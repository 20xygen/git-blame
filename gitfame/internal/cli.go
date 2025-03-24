package internal

import (
	"fmt"
	"github.com/spf13/cobra"
)

// go run gitfame/cmd/gitfame/main.go --extensions '.go,.md' --languages 'go,markdown' --exclude 'vendor/*'
var (
	rootCmd = &cobra.Command{
		Use:   "gitfame",
		Short: "Analyze git blame statistics",
		Long:  "Gitfame is a simple CLI tool to analyze git blame directory statistics.",
		Args:  cobra.NoArgs,
		Run:   cliTest,
	}
)

func cliTest(cmd *cobra.Command, args []string) {
	path, _ := cmd.Flags().GetString("repository")
	revision, _ := cmd.Flags().GetString("revision")
	orderBy, _ := cmd.Flags().GetString("order-by")
	useCommitter, _ := cmd.Flags().GetBool("use-committer")
	exts, _ := cmd.Flags().GetStringSlice("extensions")
	langs, _ := cmd.Flags().GetStringSlice("languages")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	restrict, _ := cmd.Flags().GetStringSlice("restrict-to")

	fmt.Printf("path\t%s\n", path)
	fmt.Printf("revision\t%s\n", revision)
	fmt.Printf("orderBy\t%s\n", orderBy)
	fmt.Printf("useCommitter\t%t\n", useCommitter)
	fmt.Printf("exts\t%v\n", exts)
	fmt.Printf("languages\t%v\n", langs)
	fmt.Printf("exclude\t%v\n", exclude)
	fmt.Printf("restrict\t%v\n", restrict)
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringP("repository", "r", ".", "Git repository path")
	rootCmd.Flags().StringP("revision", "R", "HEAD", "Git revision")
	rootCmd.Flags().StringP("order-by", "o", "lines", "Sort key: lines/commits/files")
	rootCmd.Flags().BoolP("use-committer", "C", false, "Use committer instead of author")
	rootCmd.Flags().StringSliceP("extensions", "e", nil, "File extensions filter (comma-separated)")
	rootCmd.Flags().StringSliceP("languages", "l", nil, "Languages filter (comma-separated)")
	rootCmd.Flags().StringSliceP("exclude", "x", nil, "Exclude glob patterns")
	rootCmd.Flags().StringSliceP("restrict-to", "t", nil, "Restrict-to glob patterns")
	rootCmd.Flags().StringP("format", "f", "tabular", "Output format")
}
