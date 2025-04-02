package files

import (
	"path/filepath"
	"strings"
)

type LangInfo struct {
	ExtToLang map[string]string
	LangToExs map[string][]string
}

type Entity interface {
	Path() string
}

type File struct {
	Name string
	Dad  *Dir
}

func (f *File) Path() string {
	if f.Dad != nil {
		return f.Dad.Path() + "/" + f.Name
	}
	return f.Name
}

func (f *File) Rel(parent string) (string, error) {
	path, err := filepath.Rel(parent, f.Path())
	if err != nil {
		return "", ErrorRelativePath{
			E: err,
		}
	}
	return path, nil
}

func (f *File) Extension() string {
	parts := strings.Split(f.Name, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

func (f *File) Language(info *LangInfo) string {
	ext := f.Extension()
	if info != nil && info.ExtToLang != nil {
		return info.ExtToLang[ext]
	}
	return ""
}

func (f *File) Lang(info *LangInfo) string {
	ext := f.Extension()
	if info != nil && info.ExtToLang != nil {
		return info.ExtToLang[ext]
	}
	return ""
}

type Dir struct {
	Name string
	Dad  *Dir
	Kids map[string]Entity // *File or *Dir
}

func (d *Dir) Path() string {
	if d.Dad != nil {
		return d.Dad.Path() + "/" + d.Name
	}
	return d.Name
}

func (d *Dir) Walk(fn func(*File) error) error {
	for _, kid := range d.Kids {
		if fl, ok := kid.(*File); ok {
			err := fn(fl)
			if err != nil {
				return err
			}
		} else if sub, ok := kid.(*Dir); ok {
			err := sub.Walk(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
