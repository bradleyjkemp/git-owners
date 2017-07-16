package git

import (
	"bytes"
	"github.com/pkg/errors"
	"os/exec"
	"os"
	"path/filepath"
	"strings"
)

// Returns relative path to the repo root
func RepoRoot() (string, error) {
	rootBytes, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	root := string(rootBytes)
	if err != nil {
		return "", errors.Wrapf(err, "not in a git repo: %s", root)
	}
	
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
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

func trimAndSplitNull(in []byte) []string {
	return strings.Split(strings.TrimSpace(string(bytes.Trim(in, "\x00"))), "\x00")
}

// func mapFilesToOwners(fileNames []string) (map[string][]string, error) {
// 	fileToOwners := make(map[string][]string)
// 	for _, fname := range fileNames {
// 		attr, err := exec.Command("git", "check-attr", "-z", "owners", fname).CombinedOutput()
// 		if err != nil {
// 			return nil, errors.Wrapf(err, "error running git merge-base %s", string(attr))
// 		}

// 		output := trimAndSplitNull(attr)
// 		if len(output) != 3 {
// 			return nil, errors.Errorf("got invalid output from git check-attr: %v", output)
// 		}

// 		ownerList := output[2]

// 		fileToOwners[fname] = strings.Split(ownerList, ",")
// 	}

// 	return fileToOwners, nil
// }
