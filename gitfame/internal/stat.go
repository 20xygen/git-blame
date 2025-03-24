package internal

import (
	"path/filepath"
)

type statFile struct {
	linesNum int
}

type statCommit struct {
	files map[string]statFile
}

type statUser struct {
	commits map[string]*statCommit
}

func (u *statUser) totalLines() int {
	var sum int
	for _, com := range u.commits {
		for _, fl := range com.files {
			sum += fl.linesNum
		}
	}
	return sum
}

func (u *statUser) totalFiles() int {
	var sum int
	for _, com := range u.commits {
		sum += len(com.files)
	}
	return sum
}

func (u *statUser) totalCommits() int {
	return len(u.commits)
}

type statVals struct {
	Commits int `json:"commits"`
	Files   int `json:"files"`
	Lines   int `json:"lines"`
}

func (u *statUser) total() (vals statVals) {
	vals.Commits = len(u.commits)
	for _, com := range u.commits {
		vals.Files += len(com.files)
		for _, fl := range com.files {
			vals.Lines += fl.linesNum
		}
	}
	return
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

		path := fl.path()

		if len(ps.exclude) > 0 && anyFn(ps.exclude, func(s string) bool {
			ok, _ := filepath.Match(s, path) // TODO: hale wrong pattern
			return ok
		}) {
			//fmt.Println("Excluded")
			return false
		}

		if len(ps.restrict) > 0 && allFn(ps.restrict, func(s string) bool {
			ok, _ := filepath.Match(s, path)
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
	for _, ln := range bo.lines {
		var name string
		if !ps.useCommitter {
			name = ln.com.meta["author"]
		} else {
			name = ln.com.meta["committer"]
		}

		usr, ok := st.users[name]
		if !ok {
			usr = &statUser{
				commits: make(map[string]*statCommit),
			}
			st.users[name] = usr
		}

		com, ok := usr.commits[ln.com.hash]
		if !ok {
			com = &statCommit{
				files: make(map[string]statFile),
			}
			usr.commits[ln.com.hash] = com
		}

		com.files[fl.path()] = statFile{com.files[fl.path()].linesNum + 1}
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

	return st, err
}
