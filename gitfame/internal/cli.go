package internal

import (
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
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
	orderBy      []string
	useCommitter bool
	extensions   []string
	languages    []string
	exclude      []string
	restrict     []string
	format       string
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
	orderBy, _ := cmd.Flags().GetStringSlice("order-by")
	useCommitter, _ := cmd.Flags().GetBool("use-committer")
	extensions, _ := cmd.Flags().GetStringSlice("extensions")
	languages, _ := cmd.Flags().GetStringSlice("languages")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	restrict, _ := cmd.Flags().GetStringSlice("restrict-to")
	format, _ := cmd.Flags().GetString("format")

	ps := &params{
		path:         path,
		revision:     revision,
		orderBy:      orderBy,
		useCommitter: useCommitter,
		extensions:   extensions,
		languages:    languages,
		exclude:      exclude,
		restrict:     restrict,
		format:       format,
	}

	fmt.Print("\nGiven parameters are:\n\n")
	fmt.Println(ps.String())

	var err error
	ps.path, err = filepath.Abs(path)
	if err != nil {
		fmt.Printf("Incorrect path: %v\n", err)
		return
	}

	fmt.Printf("\nCleaned path: %s\n\n", ps.path)

	info, err := getLangInfo()
	if err != nil {
		fmt.Printf("Error accured while loading langauges information: %v\n", err)
		return
	}

	st, err := collectStat(ps, info)
	if err != nil {
		fmt.Printf("Error accured while collecting the statistics: %v\n", err)
		return
	}

	switch format {
	case "tabular":
		fmt.Print(statTabular(st, ps.orderBy))
	case "json":
		fmt.Println(statJson(st, ps.orderBy))
	case "json-lines":
		fmt.Print(statJsonLines(st, ps.orderBy))
	case "csv":
		fmt.Print(statCSV(st, ps.orderBy))
	default:
		fmt.Print("Wrong format option.\n") // TODO: handle as error
	}
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
	rootCmd.Flags().StringP("format", "f", "tabular", "Output format")
}
