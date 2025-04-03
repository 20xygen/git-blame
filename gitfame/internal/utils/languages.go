package utils

import (
	"encoding/json"
	"gitfame/gitfame/pkg/files"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetLangInfo() (*files.LangInfo, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	var targetFile string
	if ok {
		targetFile = filepath.Dir(currentFile) + "/../../configs/language_extensions.json"
	} else {
		targetFile = "configs/language_extensions.json"
	}
	data, err := os.ReadFile(targetFile)
	if err != nil {
		return nil, ErrorConfigFile{
			E: err,
		}
	}

	var languages []struct {
		Name       string   `json:"name"`
		Extensions []string `json:"extensions"`
	}

	err = json.Unmarshal(data, &languages)
	if err != nil {
		return nil, ErrorJSONDeserialization{}
	}

	info := &files.LangInfo{
		ExtToLang: make(map[string]string),
		LangToExs: make(map[string][]string),
	}

	for _, lang := range languages {
		info.LangToExs[strings.ToLower(lang.Name)] = lang.Extensions
		for _, ext := range lang.Extensions {
			info.ExtToLang[ext] = strings.ToLower(lang.Name)
		}
	}

	return info, nil
}
