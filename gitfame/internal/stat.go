package internal

import (
	"fmt"
	"path/filepath"
)

type user struct {
	name             string
	commitToLinesNum map[string]uint64
}

func (u *user) totalLines() uint64 {
	var su uint64
	for _, v := range u.commitToLinesNum {
		su += v
	}
	return su
}

func (u *user) totalCommits() int {
	return len(u.commitToLinesNum)
}

type stat struct {
	users map[string]*user
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
			usr = &user{
				name:             name,
				commitToLinesNum: make(map[string]uint64),
			}
			st.users[name] = usr
		}

		usr.commitToLinesNum[ln.com.hash]++
	}
	return nil
}

func collectStat(ps *params, info *langInfo) (*stat, error) {
	st := &stat{
		users: make(map[string]*user),
	}

	d, err := getDirGit(ps.path, ps.revision)
	if err != nil {
		return nil, err
	}

	filter := getFileFilter(ps, info)

	err = d.walk(func(fl *file) error {
		if filter(fl) {
			fmt.Println(fl.path())
			return processFile(fl, st, ps)
		}
		return nil
	})

	return st, err
}
