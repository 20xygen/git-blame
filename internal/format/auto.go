package format

import (
	"fmt"
	"github.com/20xygen/git-blame/internal/statistics"
	"github.com/20xygen/git-blame/internal/utils"
)

func AutoFormat(st *statistics.Stat, sortKey []string, outFormat string) (string, error) {
	var tool func(*statistics.Stat, []string) (string, error)
	switch outFormat {
	case "tabular":
		tool = statTabular
	case "json":
		tool = statJSON
	case "json-lines":
		tool = statJSONLines
	case "csv":
		tool = statCSV
	case "pretty":
		tool = statPretty
	default:
		return "", utils.ErrorInvalidParameters{
			Info: fmt.Sprintf("unexpected format: %q", outFormat),
		}
	}

	return tool(st, sortKey)
}
