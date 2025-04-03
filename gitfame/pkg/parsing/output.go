package parsing

import (
	"fmt"
	"sort"
	"strings"
)

type Commit struct {
	Hash     string
	LinesNum int
	Meta     map[string]string
}

func (c *Commit) String() string {
	return fmt.Sprintf("%s\t%d\t%s\t%s", c.Hash, c.LinesNum, c.Meta["author"], c.Meta["committer"])
}

type Line struct {
	Com     *Commit
	PrevPos uint64
	CurPos  uint64
	Content string
}

type BlameOutput struct {
	Commits map[string]*Commit
	Lines   []*Line
}

func (b *BlameOutput) String() string {
	builder := strings.Builder{}
	hashToNum := make(map[string]int)

	sortedCommits := make([]*Commit, 0, len(b.Commits))
	for _, com := range b.Commits {
		sortedCommits = append(sortedCommits, com)
	}
	sort.Slice(sortedCommits, func(i, j int) bool {
		t1 := sortedCommits[i].Meta["committer-time"]
		t2 := sortedCommits[j].Meta["committer-time"]
		if t1 == "" && t2 == "" {
			return sortedCommits[i].Hash < sortedCommits[j].Hash
		}
		return t1 < t2
	})

	builder.WriteString("Num\tHash (Meta)\n")
	i := 0
	for _, com := range sortedCommits {
		hashToNum[com.Hash] = i
		builder.WriteString(fmt.Sprintf("%d\t%s\n", i, com.Hash))

		sortedKeys := make([]string, 0, len(com.Meta))
		for k := range com.Meta {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			v := com.Meta[k]
			if len(k) < 8 {
				builder.WriteString(fmt.Sprintf("\t%s\t\t%s\n", k, v))
			} else {
				builder.WriteString(fmt.Sprintf("\t%s\t%s\n", k, v))
			}
		}
		i++
	}

	builder.WriteString("\nNum\tCommit\tContent\n")

	sortedLines := make([]*Line, 0, len(b.Lines))
	sortedLines = append(sortedLines, b.Lines...)
	sort.Slice(sortedLines, func(i, j int) bool {
		return sortedLines[i].CurPos < sortedLines[j].CurPos
	})

	for _, ln := range sortedLines {
		builder.WriteString(fmt.Sprintf("%d\t%d\t%s\n", ln.CurPos, hashToNum[ln.Com.Hash], ln.Content))
	}

	return builder.String()
}
