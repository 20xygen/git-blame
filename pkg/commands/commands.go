package commands

import (
	"os/exec"
)

func commandOutput(cmd *exec.Cmd, repo string) ([]byte, error) { // TODO: move to another file
	cmd.Dir = repo
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, ErrorCommandExecution{
			C: cmd.String(),
			E: err,
		}
	}
	return out, nil
}

func GitTree(path, revision string) ([]byte, error) {
	cmd := exec.Command("git", "ls-tree", "-r", revision)
	return commandOutput(cmd, path)
}

func GitBlame(repo, path, revision string) ([]byte, error) {
	cmd := exec.Command("git", "blame", "--porcelain", revision, path)
	return commandOutput(cmd, repo)
}

func GitRevList(repo, path, revision string) ([]byte, error) {
	cmd := exec.Command("git", "rev-list", "-1", revision, "--", path)
	return commandOutput(cmd, repo)
}

func GitLog(repo, path, revision string) ([]byte, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:'%H\n%an\n%cn'", revision, "--", path)
	out, err := commandOutput(cmd, repo)
	if err != nil {
		return nil, err
	}
	return out[1 : len(out)-1], nil
}
