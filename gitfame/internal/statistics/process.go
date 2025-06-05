package statistics

import (
	"github.com/20xygen/git-blame/gitfame/internal/utils"
	"github.com/20xygen/git-blame/gitfame/pkg/files"
	"github.com/20xygen/git-blame/gitfame/pkg/parsing"
	"path/filepath"
)

func getFileFilter(ps *Params, info *files.LangInfo) func(*files.File) (bool, error) {
	return func(fl *files.File) (bool, error) {
		if len(ps.Extensions) > 0 {
			if !utils.Contains(ps.Extensions, fl.Extension()) {
				return false, nil
			}
		}

		if len(ps.Languages) > 0 {
			if !utils.Contains(ps.Languages, fl.Lang(info)) {
				return false, nil
			}
		}

		rel, err := fl.Rel(ps.Path)
		if err != nil {
			return false, err
		}

		if len(ps.Exclude) > 0 {
			for _, pat := range ps.Exclude {
				_, err := filepath.Glob(pat)
				if err != nil {
					return false, utils.ErrorInvalidPattern{
						E: err,
					}
				}
			}

			if utils.AnyFn(ps.Exclude, func(s string) bool {
				ok, _ := filepath.Match(s, rel)
				return ok
			}) {
				return false, nil
			}
		}

		if len(ps.Restrict) > 0 {
			for _, pat := range ps.Restrict {
				_, err := filepath.Glob(pat)
				if err != nil {
					return false, utils.ErrorInvalidPattern{
						E: err,
					}
				}
			}

			if utils.AllFn(ps.Restrict, func(s string) bool {
				ok, _ := filepath.Match(s, rel)
				return !ok
			}) {
				return false, nil
			}
		}

		return true, nil
	}
}

func processFile(fl *files.File, st *Stat, ps *Params) error {
	bo, err := parsing.ParseBlame(ps.Path, fl.Path(), ps.Revision)
	if err != nil {
		return err
	}
	for _, com := range bo.Commits {
		var name string
		if !ps.UseCommitter {
			name = com.Meta["author"]
		} else {
			name = com.Meta["committer"]
		}

		usr, ok := st.Users[name]
		if !ok {
			usr = &StatUser{
				Commits: make(map[string]struct{}),
				Files:   make(map[string]struct{}),
				Lines:   0,
			}
			st.Users[name] = usr
		}

		usr.Commits[com.Hash] = struct{}{}
		usr.Files[fl.Path()] = struct{}{}
		usr.Lines += com.LinesNum
	}
	return nil
}

func CollectStat(ps *Params, info *files.LangInfo) (*Stat, error) {
	st := &Stat{
		Users: make(map[string]*StatUser),
	}

	d, err := files.GetDirGit(ps.Path, ps.Revision)
	if err != nil {
		return nil, err
	}

	filter := getFileFilter(ps, info)

	err = d.Walk(func(fl *files.File) error {
		ok, errF := filter(fl)
		if errF != nil {
			return errF
		}
		if ok {
			return processFile(fl, st, ps)
		}
		return nil
	})

	return st, err
}
