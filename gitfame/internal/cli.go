package internal

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// go run gitfame/cmd/gitfame/main.go --extensions '.go,.md' --languages 'go,markdown' --exclude 'vendor/*'
var (
	rootCmd = &cobra.Command{
		Use:   "gitfame",
		Short: "Analyze git blame statistics",
		Long:  "Gitfame is a simple CLI tool to analyze git blame directory statistics.",
		Args:  cobra.NoArgs,
		Run:   command,
	}
)

type params struct {
	path         string
	revision     string
	orderBy      string
	useCommitter bool
	extensions   []string
	languages    []string
	exclude      []string
	restrict     []string
}

func (ps *params) filterLanguages(info *langInfo) error {
	var filtered []string
	flag := false
	for _, lang := range ps.languages {
		if _, ok := info.extToLang[lang]; ok {
			filtered = append(filtered, lang)
		} else {
			flag = true
		}
	}
	ps.languages = filtered
	if flag {
		return fmt.Errorf("there was an undefined language in params")
	}
	return nil
}

func (ps *params) String() string {
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "path\t\t%s\n", ps.path)
	_, _ = fmt.Fprintf(&builder, "revision\t%s\n", ps.revision)
	_, _ = fmt.Fprintf(&builder, "orderBy\t\t%s\n", ps.orderBy)
	_, _ = fmt.Fprintf(&builder, "useCommitter\t%t\n", ps.useCommitter)
	_, _ = fmt.Fprintf(&builder, "extensions\t\t%v\n", ps.extensions)
	_, _ = fmt.Fprintf(&builder, "languages\t%v\n", ps.languages)
	_, _ = fmt.Fprintf(&builder, "exclude\t\t%v\n", ps.exclude)
	_, _ = fmt.Fprintf(&builder, "restrict\t%v\n", ps.restrict)
	return builder.String()
}

func command(cmd *cobra.Command, args []string) {
	path, _ := cmd.Flags().GetString("repository")
	revision, _ := cmd.Flags().GetString("revision")
	orderBy, _ := cmd.Flags().GetString("order-by")
	useCommitter, _ := cmd.Flags().GetBool("use-committer")
	extensions, _ := cmd.Flags().GetStringSlice("extensions")
	languages, _ := cmd.Flags().GetStringSlice("languages")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	restrict, _ := cmd.Flags().GetStringSlice("restrict-to")

	ps := &params{
		path:         path,
		revision:     revision,
		orderBy:      orderBy,
		useCommitter: useCommitter,
		extensions:   extensions,
		languages:    languages,
		exclude:      exclude,
		restrict:     restrict,
	}

	fmt.Print("\nGiven parameters are:\n\n")
	fmt.Println(ps.String())

	info, err := getLangInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

	st, err := collectStat(ps, info)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Name\t\tCommits\tLines\n\n")
	for _, usr := range st.users {
		fmt.Printf("%s\t%d\t%d\n", usr.name, usr.totalCommits(), usr.totalLines())
	}
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
