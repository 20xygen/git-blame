package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"sort"
	"strings"
)

type statUnit struct {
	Name string `json:"name"`
	statVals
}

func sorted(st *stat, sortKey []string) []*statUnit {
	units := make([]*statUnit, 0, len(st.users))
	for name, user := range st.users {
		units = append(units, &statUnit{
			Name:     name,
			statVals: user.total(),
		})
	}

	sort.Slice(units, func(i, j int) bool {
		for _, key := range sortKey {
			switch strings.ToLower(key) {
			case "lines":
				if units[i].Lines != units[j].Lines {
					return units[i].Lines < units[j].Lines
				}
			case "commits":
				if units[i].Commits != units[j].Commits {
					return units[i].Commits < units[j].Commits
				}
			case "files":
				if units[i].Files != units[j].Files {
					return units[i].Files < units[j].Files
				}
			case "names":
				if units[i].Name != units[j].Name {
					return units[i].Name < units[j].Name
				}
			}
		}
		return units[i].Name < units[j].Name
	})

	return units
}

func statTabular(st *stat, sortKey []string) string {
	units := sorted(st, sortKey)

	builder := strings.Builder{}

	t := table.NewWriter()
	t.SetOutputMirror(&builder)
	t.AppendHeader(table.Row{"Name", "Commits", "Files", "Lines"})
	rows := make([]table.Row, 0, len(st.users))
	for _, unit := range units {
		rows = append(rows, table.Row{unit.Name, unit.Commits, unit.Files, unit.Lines})
	}
	t.AppendRows(rows)
	t.AppendSeparator()
	t.Render()

	return builder.String()
}

func statCSV(st *stat, sortKey []string) string {
	units := sorted(st, sortKey)
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	err := writer.Write([]string{"Name", "Lines", "Commits", "Files"})
	if err != nil {
		return "" // TODO
	}

	for _, u := range units {
		err = writer.Write([]string{
			u.Name,
			fmt.Sprintf("%d", u.Lines),
			fmt.Sprintf("%d", u.Commits),
			fmt.Sprintf("%d", u.Files),
		})
		if err != nil {
			return ""
		}
	}

	writer.Flush()

	return builder.String()
}

func statJson(st *stat, sortKey []string) string {
	units := sorted(st, sortKey)
	jsonData, err := json.MarshalIndent(units, "", "  ")
	if err != nil {
		return "" // TODO
	}
	return string(jsonData)
}

func statJsonLines(st *stat, sortKey []string) string {
	units := sorted(st, sortKey)
	var builder strings.Builder

	for _, unit := range units {
		jsonData, err := json.Marshal(unit)
		if err != nil {
			continue
		}
		builder.Write(jsonData)
		builder.WriteString("\n")
	}

	return builder.String()
}
