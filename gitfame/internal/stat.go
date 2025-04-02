package internal

import (
	"path/filepath"
	"strconv"
	"strings"
)

type statUser struct {
	commits map[string]struct{}
	files   map[string]struct{}
	lines   int
}

type statVals struct {
	Commits int `json:"commits"`
	Files   int `json:"files"`
	Lines   int `json:"lines"`
}

func (su *statUser) String() string {
	mx := 0

	var commits []string
	for hash, _ := range su.commits {
		commits = append(commits, hash)
	}
	if len(commits) > mx {
		mx = len(commits)
	}

	var files []string
	for path, _ := range su.files {
		files = append(files, path)
	}
	if len(files) > mx {
		mx = len(files)
	}

	var builder strings.Builder
	builder.WriteString(strconv.Itoa(su.lines))
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

func (su *statUser) total() statVals {
	return statVals{
		Commits: len(su.commits),
		Files:   len(su.files),
		Lines:   su.lines,
	}
}

type stat struct {
	users map[string]*statUser
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func anyFn(list []string, fn func(string) bool) bool {
	for _, item := range list {
		if fn(item) {
			return true
		}
	}
	return false
}

func allFn(list []string, fn func(string) bool) bool {
	for _, item := range list {
		if !fn(item) {
			return false
		}
	}
	return true
}

func getFileFilter(ps *params, info *langInfo) func(*file) bool {
	return func(fl *file) bool {
		if len(ps.extensions) > 0 {
			if !contains(ps.extensions, fl.extension()) {
				//fmt.Println("Bad extension")
				return false
			}
		}

		if len(ps.languages) > 0 {
			if !contains(ps.languages, fl.lang(info)) {
				//fmt.Println("Bad language")
				return false
			}
		}

		rel, _ := fl.rel(ps.path)

		if len(ps.exclude) > 0 && anyFn(ps.exclude, func(s string) bool {
			ok, _ := filepath.Match(s, rel)
			//slog.Info("Exclude match", slog.String("file", fl.name), slog.String("pattern", s), slog.Bool("result", ok))
			return ok
		}) {
			//fmt.Println("Excluded")
			return false
		}

		if len(ps.restrict) > 0 && allFn(ps.restrict, func(s string) bool {
			ok, _ := filepath.Match(s, rel)
			return !ok
		}) {
			//fmt.Println("Restricted")
			return false
		}

		return true
	}
}

func processFile(fl *file, st *stat, ps *params) error {
	bo, err := parseBlame(ps.path, fl.path(), ps.revision)
	if err != nil {
		//fmt.Printf("blame parsing error on the %s: %v\n", fl.path(), err)
		return err
	}
	for _, com := range bo.commits {
		var name string
		if !ps.useCommitter {
			name = com.meta["author"]
		} else {
			name = com.meta["committer"]
		}

		usr, ok := st.users[name]
		if !ok {
			usr = &statUser{
				commits: make(map[string]struct{}),
				files:   make(map[string]struct{}),
				lines:   0,
			}
			st.users[name] = usr
		}

		//fmt.Fprintf(os.Stderr, "Hash processed: %s\n", com.hash)

		usr.commits[com.hash] = struct{}{}
		usr.files[fl.path()] = struct{}{}
		usr.lines += com.linesNum
	}
	return nil
}

func collectStat(ps *params, info *langInfo) (*stat, error) {
	st := &stat{
		users: make(map[string]*statUser),
	}

	d, err := getDirGit(ps.path, ps.revision)
	if err != nil {
		return nil, err
	}

	filter := getFileFilter(ps, info)

	err = d.walk(func(fl *file) error {
		if filter(fl) {
			return processFile(fl, st, ps)
		}
		return nil
	})

	//for name, u := range st.users {
	//	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", name, u.String())
	//}

	return st, err
}
