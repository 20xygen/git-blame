package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func commandOutput(cmd *exec.Cmd, repo string) ([]byte, error) { // TODO: move to another file
	cmd.Dir = repo
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func blame(repo, path, revision string) ([]byte, error) {
	cmd := exec.Command("git", "blame", "--porcelain", revision, path)
	return commandOutput(cmd, repo)
}

func revList(repo, path, revision string) ([]byte, error) {
	cmd := exec.Command("git", "rev-list", "-1", revision, "--", path)
	return commandOutput(cmd, repo)
}

func gitLog(repo, path, revision string) ([]byte, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:'%H\n%an\n%cn'", revision, "--", path)
	out, err := commandOutput(cmd, repo)
	if err != nil {
		return nil, err
	}
	return out[1 : len(out)-1], nil
}

type commit struct {
	hash     string
	linesNum int
	meta     map[string]string
}

func (c *commit) String() string {
	return fmt.Sprintf("%s\t%d\t%s\t%s", c.hash, c.linesNum, c.meta["author"], c.meta["committer"])
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

func parseBlame(repo, path, revision string) (*blameOutput, error) {
	out, err := blame(repo, path, revision)
	if err != nil {
		fmt.Println(repo, path)
		return nil, err
	}

	bo := blameOutput{
		commits: make(map[string]*commit),
		lines:   make([]*line, 0),
	}

	empty := true
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		empty = false
		firstLine := scanner.Text()

		if firstLine == " <nil>" {
			continue
		}

		parts := strings.Split(firstLine, " ")
		if len(parts) < 3 || len(parts) > 4 {
			return nil, fmt.Errorf("blame: invalid line: %q", firstLine)
		}
		hash := parts[0]
		//fmt.Fprintf(os.Stderr, "Hash found: %s\n", hash)
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

	if empty {
		log, err := gitLog(repo, path, revision)
		if err != nil {
			return nil, err
		}

		slog.Info("Found empty file", slog.String("file", path), slog.String("git log", string(log)))

		parts := strings.Split(string(log), "\n")
		if len(parts) < 3 {
			return nil, fmt.Errorf("git log: invalid output %q", string(log))
		}
		hash := parts[0]
		author := parts[1]
		committer := parts[2]
		//fmt.Fprintf(os.Stderr, "Hash found: %s\n", hash)

		com := commit{
			hash:     hash,
			linesNum: 0,
			meta:     make(map[string]string),
		}

		com.meta["author"] = author
		com.meta["committer"] = committer
		bo.commits[hash] = &com

	} else {
		for _, ln := range bo.lines {
			ln.com.linesNum++
		}
	}

	//var builder strings.Builder
	//builder.WriteString(fmt.Sprintf("File %s has %d commits:\n", path, len(bo.commits)))
	//for _, com := range bo.commits {
	//	builder.WriteString("\t")
	//	builder.WriteString(com.String())
	//	builder.WriteString("\n")
	//}
	//builder.WriteString("\n")
	//_, _ = fmt.Fprintf(os.Stderr, builder.String())

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
