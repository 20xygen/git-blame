package statistics

import (
	"fmt"
	"github.com/20xygen/git-blame/internal/utils"
	"github.com/20xygen/git-blame/pkg/files"
	"github.com/spf13/cobra"
	"strings"
)

type Params struct {
	Path         string
	Revision     string
	OrderBy      []string
	UseCommitter bool
	Extensions   []string
	Languages    []string
	Exclude      []string
	Restrict     []string
	Format       string
}

func (ps *Params) FilterLanguages(info *files.LangInfo) error {
	var filtered []string
	flag := false
	for _, lang := range ps.Languages {
		if _, ok := info.ExtToLang[lang]; ok {
			filtered = append(filtered, lang)
		} else {
			flag = true
		}
	}
	ps.Languages = filtered
	if flag {
		return utils.ErrorUndefinedLanguage{}
	}
	return nil
}

func (ps *Params) String() string {
	var builder strings.Builder
	_, _ = fmt.Fprintf(&builder, "path\t\t%s\n", ps.Path)
	_, _ = fmt.Fprintf(&builder, "revision\t%s\n", ps.Revision)
	_, _ = fmt.Fprintf(&builder, "orderBy\t\t%s\n", ps.OrderBy)
	_, _ = fmt.Fprintf(&builder, "useCommitter\t%t\n", ps.UseCommitter)
	_, _ = fmt.Fprintf(&builder, "extensions\t\t%v\n", ps.Extensions)
	_, _ = fmt.Fprintf(&builder, "languages\t%v\n", ps.Languages)
	_, _ = fmt.Fprintf(&builder, "exclude\t\t%v\n", ps.Exclude)
	_, _ = fmt.Fprintf(&builder, "restrict\t%v\n", ps.Restrict)
	return builder.String()
}

func GetParams(cmd cobra.Command) (*Params, error) {
	path, e1 := cmd.Flags().GetString("repository")
	revision, e2 := cmd.Flags().GetString("revision")
	orderBy, e3 := cmd.Flags().GetStringSlice("order-by")
	useCommitter, e4 := cmd.Flags().GetBool("use-committer")
	extensions, e5 := cmd.Flags().GetStringSlice("extensions")
	languages, e6 := cmd.Flags().GetStringSlice("languages")
	exclude, e7 := cmd.Flags().GetStringSlice("exclude")
	restrict, e8 := cmd.Flags().GetStringSlice("restrict-to")
	formatArg, e9 := cmd.Flags().GetString("format")

	if utils.AnyError(e1, e2, e3, e4, e5, e6, e7, e8, e9) {
		return nil, utils.ErrorInvalidParameters{
			Info: "unexpected error",
		}
	}

	return &Params{
		Path:         path,
		Revision:     revision,
		OrderBy:      orderBy,
		UseCommitter: useCommitter,
		Extensions:   extensions,
		Languages:    languages,
		Exclude:      exclude,
		Restrict:     restrict,
		Format:       formatArg,
	}, nil
}
