package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type langInfo struct {
	extToLang map[string]string
}

func getLangInfo() (*langInfo, error) {
	data, err := os.ReadFile("gitfame/configs/language_extensions.json")
	if err != nil {
		return nil, err
	}

	var languages []struct {
		Name       string   `json:"name"`
		Extensions []string `json:"extensions"`
	}

	err = json.Unmarshal(data, &languages)
	if err != nil {
		return nil, err
	}

	info := &langInfo{
		extToLang: make(map[string]string),
	}

	for _, lang := range languages {
		for _, ext := range lang.Extensions {
			info.extToLang[ext] = lang.Name
		}
	}

	return info, nil
}

type entity interface {
	path() string
}

type file struct {
	name string
	dad  *dir
}

func (f *file) path() string {
	if f.dad != nil {
		return f.dad.path() + "/" + f.name
	}
	return f.name
}

func (f *file) extension() string {
	parts := strings.Split(f.name, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

func (f *file) lang(info *langInfo) string {
	ext := f.extension()
	if info != nil && info.extToLang != nil {
		return info.extToLang[ext]
	}
	return ""
}

type dir struct {
	name string
	dad  *dir
	kids map[string]entity // *file or *dir
}

func (d *dir) path() string {
	if d.dad != nil {
		return d.dad.path() + "/" + d.name
	}
	return d.name
}

type fileInfo struct {
	fl   *file
	ext  string
	lang string
}

func (d *dir) files() []fileInfo {
	info, err := getLangInfo()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return d.filesHelper(nil, info)
}

func (d *dir) filesHelper(list []fileInfo, info *langInfo) []fileInfo {
	for _, kid := range d.kids {
		if fl, ok := kid.(*file); ok {
			list = append(list, fileInfo{
				fl:   fl,
				ext:  fl.extension(),
				lang: fl.lang(info),
			})
		} else if sub, ok := kid.(*dir); ok {
			list = sub.filesHelper(list, info)
		}
	}
	return list
}

func (d *dir) asList() string {
	var builder strings.Builder
	d.asListHelper(&builder)
	return builder.String()
}

func (d *dir) asListHelper(builder *strings.Builder) {
	var keys []string
	for k := range d.kids {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	builder.WriteString(d.path())
	builder.WriteString("\n")

	for _, k := range keys {
		kid := d.kids[k]
		if sub, ok := kid.(*dir); ok {
			sub.asListHelper(builder)
		}
	}
	for _, k := range keys {
		kid := d.kids[k]
		if fl, ok := kid.(*file); ok {
			builder.WriteString(fl.path())
			builder.WriteString("\n")
		}
	}

}

func getDir(path string) (*dir, error) {
	root := dir{
		name: path,
		kids: make(map[string]entity),
	}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(path, filePath)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		components := strings.Split(relPath, string(filepath.Separator))

		currentDir := &root
		for i := 0; i < len(components)-1; i++ {
			dirName := components[i]
			if _, exists := currentDir.kids[dirName]; !exists {
				newDir := dir{
					name: dirName,
					dad:  currentDir,
					kids: make(map[string]entity),
				}
				currentDir.kids[dirName] = &newDir
			}
			currentDir = currentDir.kids[dirName].(*dir)
		}

		name := components[len(components)-1]
		if info.IsDir() {
			newDir := dir{
				name: name,
				dad:  currentDir,
				kids: make(map[string]entity),
			}
			currentDir.kids[name] = &newDir
		} else {
			newFile := file{
				name: name,
				dad:  currentDir,
			}
			currentDir.kids[name] = &newFile
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &root, nil
}
