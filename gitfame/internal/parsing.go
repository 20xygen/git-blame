package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func blame(path string) ([]byte, error) {
	out, err := exec.Command("git", "blame", "--porcelain", path).Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

type commit struct {
	hash     string
	linesNum uint64
	meta     map[string]string
}

type line struct {
	com     *commit
	prevPos uint64
	curPos  uint64
	content string
}

type blameOutput struct {
	commits map[string]*commit
	lines   []*line
}

func parseBlame(path string) (*blameOutput, error) {
	out, err := blame(path)
	if err != nil {
		return nil, err
	}

	bo := blameOutput{
		commits: make(map[string]*commit),
		lines:   make([]*line, 0),
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		firstLine := scanner.Text()

		if firstLine == " <nil>" {
			continue
		}

		parts := strings.Split(firstLine, " ")
		if len(parts) < 3 || len(parts) > 4 {
			return nil, fmt.Errorf("blame: invalid line: %q", firstLine)
		}
		hash := parts[0]
		prev, err1 := strconv.ParseUint(parts[1], 10, 64)
		cur, err2 := strconv.ParseUint(parts[2], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("blame: invalid line: %q", firstLine)
		}

		com, ok := bo.commits[hash]
		if !ok {
			com = &commit{
				hash: hash,
				meta: make(map[string]string),
			}
			bo.commits[hash] = com
		}

		if !scanner.Scan() {
			return nil, fmt.Errorf("blame: invalid block after %v (empty block)", parts)
		}
		nextLine := scanner.Text()
		for nextLine[0] != '\t' {
			if nextLine != "boundary" {
				params := strings.SplitN(nextLine, " ", 2)
				if len(params) != 2 {
					return nil, fmt.Errorf("blame: invalid block line %q", nextLine)
				}
				com.meta[params[0]] = params[1]
			}

			if !scanner.Scan() {
				return nil, fmt.Errorf("blame: invalid block after %v (block with no content)", parts)
			}
			nextLine = scanner.Text()
		}

		bo.lines = append(bo.lines, &line{
			com:     com,
			prevPos: prev,
			curPos:  cur,
			content: nextLine[1:],
		})
	}

	return &bo, nil
}

func (b *blameOutput) String() string {
	builder := strings.Builder{}
	hashToNum := make(map[string]int)

	sortedCommits := make([]*commit, 0, len(b.commits))
	for _, com := range b.commits {
		sortedCommits = append(sortedCommits, com)
	}
	sort.Slice(sortedCommits, func(i, j int) bool {
		t1 := sortedCommits[i].meta["committer-time"]
		t2 := sortedCommits[j].meta["committer-time"]
		if t1 == "" && t2 == "" {
			return sortedCommits[i].hash < sortedCommits[j].hash
		}
		return t1 < t2
	})

	builder.WriteString("Num\tHash (Meta)\n")
	i := 0
	for _, com := range sortedCommits {
		hashToNum[com.hash] = i
		builder.WriteString(fmt.Sprintf("%d\t%s\n", i, com.hash))

		sortedKeys := make([]string, 0, len(com.meta))
		for k := range com.meta {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			v := com.meta[k]
			if len(k) < 8 {
				builder.WriteString(fmt.Sprintf("\t%s\t\t%s\n", k, v))
			} else {
				builder.WriteString(fmt.Sprintf("\t%s\t%s\n", k, v))
			}
		}
		i++
	}

	builder.WriteString("\nNum\tCommit\tContent\n")

	sortedLines := make([]*line, 0, len(b.lines))
	for _, ln := range b.lines {
		sortedLines = append(sortedLines, ln)
	}
	sort.Slice(sortedLines, func(i, j int) bool {
		return sortedLines[i].curPos < sortedLines[j].curPos
	})

	for _, ln := range sortedLines {
		builder.WriteString(fmt.Sprintf("%d\t%d\t%s\n", ln.curPos, hashToNum[ln.com.hash], ln.content))
	}

	return builder.String()
}
