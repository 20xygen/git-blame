package format

import (
	"fmt"
	"gitfame/gitfame/internal/statistics"
	"gitfame/gitfame/internal/utils"
)

func AutoFormat(st *statistics.Stat, sortKey []string, outFormat string) (string, error) {
	var tool func(*statistics.Stat, []string) (string, error)
	switch outFormat {
	case "tabular":
		tool = statTabular
	case "json":
		tool = statJson
	case "json-lines":
		tool = statJsonLines
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
