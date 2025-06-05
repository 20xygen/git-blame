package files

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/20xygen/git-blame/gitfame/pkg/commands"
)

type FileInfo struct {
	fl   *File
	ext  string
	lang string
}

func (d *Dir) Files(info *LangInfo) []FileInfo {
	return d.filesHelper(nil, info)
}

func (d *Dir) filesHelper(list []FileInfo, info *LangInfo) []FileInfo {
	for _, kid := range d.Kids {
		if fl, ok := kid.(*File); ok {
			list = append(list, FileInfo{
				fl:   fl,
				ext:  fl.Extension(),
				lang: fl.Language(info),
			})
		} else if sub, ok := kid.(*Dir); ok {
			list = sub.filesHelper(list, info)
		}
	}
	return list
}

func (d *Dir) AsList() string {
	var builder strings.Builder
	d.asListHelper(&builder)
	return builder.String()
}

func (d *Dir) asListHelper(builder *strings.Builder) {
	var keys []string
	for k := range d.Kids {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	builder.WriteString(d.Path())
	builder.WriteString("\n")

	for _, k := range keys {
		kid := d.Kids[k]
		if sub, ok := kid.(*Dir); ok {
			sub.asListHelper(builder)
		} else if fl, ok := kid.(*File); ok {
			builder.WriteString(fl.Path())
			builder.WriteString("\n")
		}
	}
}

func gitTreePaths(path, revision string) ([]string, error) {
	out, err := commands.GitTree(path, revision)
	if err != nil {
		return nil, err
	}

	var paths []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		ln := scanner.Text()
		parts := strings.Split(ln, "\t")
		if len(parts) != 2 {
			return nil, commands.ErrorInvalidGitTreeOutput{}
		}
		paths = append(paths, parts[1])
	}

	return paths, nil
}

func getDirPaths(rootPath string, paths []string) *Dir {
	root := &Dir{
		Name: rootPath,
		Kids: make(map[string]Entity),
	}

	for _, path := range paths {
		components := strings.Split(path, string(filepath.Separator))

		currentDir := root
		for i := 0; i < len(components)-1; i++ {
			dirName := components[i]
			if _, exists := currentDir.Kids[dirName]; !exists {
				newDir := Dir{
					Name: dirName,
					Dad:  currentDir,
					Kids: make(map[string]Entity),
				}
				currentDir.Kids[dirName] = &newDir
			}
			currentDir = currentDir.Kids[dirName].(*Dir)
		}

		name := components[len(components)-1]
		newFile := File{
			Name: name,
			Dad:  currentDir,
		}
		currentDir.Kids[name] = &newFile
	}

	return root
}

func GetDir(path string) (*Dir, error) {
	var paths *[]string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(path, filePath)
		if err != nil {
			return ErrorWalk{
				Err: err,
			}
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

func GetDirGit(path, revision string) (*Dir, error) {
	paths, err := gitTreePaths(path, revision)
	if err != nil {
		return nil, err
	}
	return getDirPaths(path, paths), nil
}
