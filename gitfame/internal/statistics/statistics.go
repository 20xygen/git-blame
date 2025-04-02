package statistics

import (
	"strconv"
	"strings"
)

type StatUser struct {
	Commits map[string]struct{}
	Files   map[string]struct{}
	Lines   int
}

type StatVals struct {
	Commits int `json:"commits"`
	Files   int `json:"files"`
	Lines   int `json:"lines"`
}

type Stat struct {
	Users map[string]*StatUser
}

func (su *StatUser) String() string {
	mx := 0

	var commits []string
	for hash, _ := range su.Commits {
		commits = append(commits, hash)
	}
	if len(commits) > mx {
		mx = len(commits)
	}

	var files []string
	for path, _ := range su.Files {
		files = append(files, path)
	}
	if len(files) > mx {
		mx = len(files)
	}

	var builder strings.Builder
	builder.WriteString(strconv.Itoa(su.Lines))
	builder.WriteString(" lines\n#\tCommits")
	builder.WriteString(strings.Repeat(" ", 33))
	builder.WriteString("\tFiles\n")

	for i := range mx {
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString("\t")

		if len(commits) > i {
			builder.WriteString(commits[i])
		} else {
			builder.WriteString(strings.Repeat(" ", 40))
		}
		builder.WriteString("\t")

		if len(files) > i {
			builder.WriteString(files[i])
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func (su *StatUser) Total() StatVals {
	return StatVals{
		Commits: len(su.Commits),
		Files:   len(su.Files),
		Lines:   su.Lines,
	}
}
