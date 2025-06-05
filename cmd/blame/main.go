//go:build !solution

package main

import (
	"github.com/20xygen/git-blame/internal/cli"
)

func main() {
	_ = cli.Execute()
}
