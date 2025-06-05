package parsing

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/20xygen/git-blame/gitfame/pkg/commands"
)

func parseEmpty(repo, path, revision string, bo *BlameOutput) error {
	log, err := commands.GitLog(repo, path, revision)
	if err != nil {
		return err
	}

	parts := strings.Split(string(log), "\n")
	if len(parts) < 3 {
		return commands.ErrorInvalidGitLogOutput{}
	}
	hash := parts[0]
	author := parts[1]
	committer := parts[2]

	com := Commit{
		Hash:     hash,
		LinesNum: 0,
		Meta:     make(map[string]string),
	}

	com.Meta["author"] = author
	com.Meta["committer"] = committer
	bo.Commits[hash] = &com

	return nil
}

func ParseBlame(repo, path, revision string) (*BlameOutput, error) {
	out, err := commands.GitBlame(repo, path, revision)
	if err != nil {
		fmt.Println(repo, path)
		return nil, err
	}

	bo := BlameOutput{
		Commits: make(map[string]*Commit),
		Lines:   make([]*Line, 0),
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
			return nil, commands.ErrorInvalidGitBlameOutput{
				Info: fmt.Sprintf("invalid line: %q", firstLine),
			}
		}
		hash := parts[0]
		prev, err1 := strconv.ParseUint(parts[1], 10, 64)
		cur, err2 := strconv.ParseUint(parts[2], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, commands.ErrorInvalidGitBlameOutput{
				Info: fmt.Sprintf("invalid line: %q", firstLine),
			}
		}

		com, ok := bo.Commits[hash]
		if !ok {
			com = &Commit{
				Hash: hash,
				Meta: make(map[string]string),
			}
			bo.Commits[hash] = com
		}

		if !scanner.Scan() {
			return nil, commands.ErrorInvalidGitBlameOutput{
				Info: fmt.Sprintf("invalid block after %v (empty block)", parts),
			}
		}
		nextLine := scanner.Text()
		for nextLine[0] != '\t' {
			if nextLine != "boundary" {
				params := strings.SplitN(nextLine, " ", 2)
				if len(params) != 2 {
					return nil, commands.ErrorInvalidGitBlameOutput{
						Info: fmt.Sprintf("invalid block line %q", nextLine),
					}
				}
				com.Meta[params[0]] = params[1]
			}

			if !scanner.Scan() {
				return nil, commands.ErrorInvalidGitBlameOutput{
					Info: fmt.Sprintf("invalid block after %v (block with no content)", parts),
				}
			}
			nextLine = scanner.Text()
		}

		bo.Lines = append(bo.Lines, &Line{
			Com:     com,
			PrevPos: prev,
			CurPos:  cur,
			Content: nextLine[1:],
		})
	}

	if empty {
		err := parseEmpty(repo, path, revision, &bo)
		if err != nil {
			return nil, err
		}

	} else {
		for _, ln := range bo.Lines {
			ln.Com.LinesNum++
		}
	}

	return &bo, nil
}
