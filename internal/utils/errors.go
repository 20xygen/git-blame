package utils

import "fmt"

const (
	_ = iota
	CodeParametersParsing
	CodeAbsolutePath
	CodeLanguageInfo
	CodeFormat
)

type ErrorUndefinedLanguage struct{}

func (e ErrorUndefinedLanguage) Error() string {
	return "there was an undefined language in parameters"
}

type ErrorInvalidParameters struct {
	Info string
}

func (e ErrorInvalidParameters) Error() string {
	if e.Info != "" {
		return fmt.Sprintf("invalid parameters (%s)", e.Info)
	}
	return "invalid parameters"
}

type ErrorJSONSerialization struct{}

func (e ErrorJSONSerialization) Error() string { return "json serialization error" }

type ErrorJSONDeserialization struct{}

func (e ErrorJSONDeserialization) Error() string { return "json deserialization error" }

type ErrorConfigFile struct {
	E error
}

func (e ErrorConfigFile) Error() string {
	return fmt.Sprintf("config file not found (%v)", e.E)
}

type ErrorInvalidPattern struct {
	E error
}

func (e ErrorInvalidPattern) Error() string {
	return fmt.Sprintf("invalid glob pattern (%v)", e.E)
}
