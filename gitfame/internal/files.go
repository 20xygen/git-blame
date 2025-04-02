package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type langInfo struct {
	extToLang map[string]string
	langToExs map[string][]string
}

func getLangInfo() (*langInfo, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	var targetFile string
	if ok {
		targetFile = filepath.Dir(currentFile) + "/../configs/language_extensions.json"
	} else {
		targetFile = "configs/language_extensions.json"
	}
	data, err := os.ReadFile(targetFile)
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
		langToExs: make(map[string][]string),
	}

	for _, lang := range languages {
		info.langToExs[strings.ToLower(lang.Name)] = lang.Extensions
		for _, ext := range lang.Extensions {
			info.extToLang[ext] = strings.ToLower(lang.Name)
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

func (f *file) rel(parent string) (string, error) {
	path, err := filepath.Rel(parent, f.path())
	if err != nil {
		return "", err
	}
	return path, nil
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
		} else if fl, ok := kid.(*file); ok {
			builder.WriteString(fl.path())
			builder.WriteString("\n")
		}
	}
}

func gitTree(path, revision string) ([]byte, error) {
	cmd := exec.Command("git", "ls-tree", "-r", revision)
	cmd.Dir = path

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git ls-tree failed: %s", err)
	}
	return out, nil
}

func gitTreePaths(path, revision string) ([]string, error) {
	out, err := gitTree(path, revision)
	if err != nil {
		return nil, err
	}

	var paths []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		ln := scanner.Text()
		parts := strings.Split(ln, "\t")
		if len(parts) != 2 {
			return nil, fmt.Errorf("wrong output format (git ls-tree):%s", ln)
		}
		paths = append(paths, parts[1])
	}

	return paths, nil
}

func getDirPaths(rootPath string, paths []string) *dir {
	root := &dir{
		name: rootPath,
		kids: make(map[string]entity),
	}

	for _, path := range paths {
		components := strings.Split(path, string(filepath.Separator))

		currentDir := root
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
		newFile := file{
			name: name,
			dad:  currentDir,
		}
		currentDir.kids[name] = &newFile
	}

	return root
}

func getDir(path string) (*dir, error) {
	var paths *[]string

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

		*paths = append(*paths, relPath)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return getDirPaths(path, *paths), nil
}

func getDirGit(path, revision string) (*dir, error) {
	paths, err := gitTreePaths(path, revision)
	if err != nil {
		return nil, err
	}
	return getDirPaths(path, paths), nil
}

func (d *dir) walk(fn func(*file) error) error {
	for _, kid := range d.kids {
		if fl, ok := kid.(*file); ok {
			err := fn(fl)
			if err != nil {
				return err
			}
		} else if sub, ok := kid.(*dir); ok {
			err := sub.walk(fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
