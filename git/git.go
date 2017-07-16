package git

import (
	"bytes"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Returns relative path to the repo root
func RepoRoot() (string, error) {
	rootBytes, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	root := strings.TrimSpace(string(rootBytes))
	if err != nil {
		return "", errors.Wrapf(err, "not in a git repo: %s", root)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	if currentDir == root {
		return ".", nil
	}

	path, err := filepath.Rel(currentDir, root)
	if err != nil {
		return "", errors.Wrap(err, "failed to construct relative path")
	}

	return path, nil
}

func FindBaseCommit(baseBranch string) (string, error) {
	baseCommit, err := exec.Command("git", "merge-base", baseBranch, "HEAD").CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "error running git merge-base: %s", string(baseCommit))
	}
	return strings.TrimSpace(string(baseCommit)), nil
}

func FindChangedFiles(startCommit string) ([]string, error) {
	changedFiles, err := exec.Command("git", "diff", "-z", "--name-only", startCommit+"..HEAD").CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "error running git diff: %s", string(changedFiles))
	}

	files := trimAndSplitNull(changedFiles)
	for i, fname := range files {
		files[i] = strings.TrimSpace(fname)
	}

	return files, nil
}

func FileContentsAtCommit(path, commit string) string {
	contents, err := exec.Command("git", "show", "commit:"+path).CombinedOutput()
	if err != nil {
		return ""
	}

	return string(contents)
}

func trimAndSplitNull(in []byte) []string {
	return strings.Split(strings.TrimSpace(string(bytes.Trim(in, "\x00"))), "\x00")
}
