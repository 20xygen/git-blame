package format

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gitfame/gitfame/internal/statistics"
	"gitfame/gitfame/internal/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"sort"
	"strings"
	"text/tabwriter"
)

type statUnit struct {
	Name string `json:"name"`
	statistics.StatVals
}

func validateSortKey(sortKey []string) error {
	valid := []string{"lines", "commits", "files", "names"}
	for _, key := range sortKey {
		if !utils.Contains(valid, key) {
			return utils.ErrorInvalidParameters{
				Info: fmt.Sprintf("unexpected sort key: %s", key),
			}
		}
	}
	return nil
}

func sorted(st *statistics.Stat, sortKey []string) ([]*statUnit, error) {
	units := make([]*statUnit, 0, len(st.Users))
	for name, user := range st.Users {
		units = append(units, &statUnit{
			Name:     name,
			StatVals: user.Total(),
		})
	}

	if err := validateSortKey(sortKey); err != nil {
		return nil, err
	}

	fullSortKey := append([]string{}, sortKey...)
	fullSortKey = append(fullSortKey, "lines")
	fullSortKey = append(fullSortKey, "commits")
	fullSortKey = append(fullSortKey, "files")

	sort.Slice(units, func(i, j int) bool {
		for _, key := range fullSortKey {
			switch strings.ToLower(key) {
			case "lines":
				if units[i].Lines != units[j].Lines {
					return units[i].Lines > units[j].Lines
				}
			case "commits":
				if units[i].Commits != units[j].Commits {
					return units[i].Commits > units[j].Commits
				}
			case "files":
				if units[i].Files != units[j].Files {
					return units[i].Files > units[j].Files
				}
			case "names":
				if units[i].Name != units[j].Name {
					return units[i].Name < units[j].Name
				}
			}
		}
		return units[i].Name < units[j].Name
	})

	return units, nil
}

func statTabular(st *statistics.Stat, sortKey []string) (string, error) {
	units, err := sorted(st, sortKey)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	writer := tabwriter.NewWriter(&builder, 0, 0, 1, ' ', 0)

	_, _ = fmt.Fprintln(writer, "Name\tLines\tCommits\tFiles")
	for _, unit := range units {
		_, _ = fmt.Fprintf(writer, "%s\t%d\t%d\t%d\n", unit.Name, unit.Lines, unit.Commits, unit.Files)
	}

	_ = writer.Flush()
	return builder.String(), nil
}

func statPretty(st *statistics.Stat, sortKey []string) (string, error) {
	units, err := sorted(st, sortKey)
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	t := table.NewWriter()
	t.SetOutputMirror(&builder)
	t.AppendHeader(table.Row{"Name", "Commits", "Files", "Lines"})
	rows := make([]table.Row, 0, len(st.Users))
	for _, unit := range units {
		rows = append(rows, table.Row{unit.Name, unit.Commits, unit.Files, unit.Lines})
	}
	t.AppendRows(rows)
	t.AppendSeparator()
	t.Render()

	return builder.String(), nil
}

func statCSV(st *statistics.Stat, sortKey []string) (string, error) {
	units, err := sorted(st, sortKey)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	err = writer.Write([]string{"Name", "Lines", "Commits", "Files"})
	if err != nil {
		return "", err
	}

	for _, u := range units {
		err = writer.Write([]string{
			u.Name,
			fmt.Sprintf("%d", u.Lines),
			fmt.Sprintf("%d", u.Commits),
			fmt.Sprintf("%d", u.Files),
		})
		if err != nil {
			return "", err
		}
	}

	writer.Flush()

	return builder.String(), nil
}

func statJSON(st *statistics.Stat, sortKey []string) (string, error) {
	units, err := sorted(st, sortKey)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(units, "", "  ")
	if err != nil {
		return "", utils.ErrorJSONSerialization{}
	}
	return string(jsonData), nil
}

func statJSONLines(st *statistics.Stat, sortKey []string) (string, error) {
	units, err := sorted(st, sortKey)
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	for _, unit := range units {
		jsonData, err := json.Marshal(unit)
		if err != nil {
			return "", utils.ErrorJSONSerialization{}
		}
		builder.Write(jsonData)
		builder.WriteString("\n")
	}

	return builder.String(), nil
}
