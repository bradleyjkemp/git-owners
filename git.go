package main

import (
	"bytes"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

func trimAndSplitNull(in []byte) []string {
	return strings.Split(strings.TrimSpace(string(bytes.Trim(in, "\x00"))), "\x00")
}

func findBaseCommit(baseBranch string) (string, error) {
	baseCommit, err := exec.Command("git", "merge-base", baseBranch, "HEAD").CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "error running git merge-base: %s", string(baseCommit))
	}
	return strings.TrimSpace(string(baseCommit)), nil
}

func findChangedFiles(startCommit string) ([]string, error) {
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

func mapFilesToOwners(fileNames []string) (map[string][]string, error) {
	fileToOwners := make(map[string][]string)
	for _, fname := range fileNames {
		attr, err := exec.Command("git", "check-attr", "-z", "owners", fname).CombinedOutput()
		if err != nil {
			return nil, errors.Wrapf(err, "error running git merge-base %s", string(attr))
		}

		output := trimAndSplitNull(attr)
		if len(output) != 3 {
			return nil, errors.Errorf("got invalid output from git check-attr: %v", output)
		}

		ownerList := output[2]

		fileToOwners[fname] = strings.Split(ownerList, ",")
	}

	return fileToOwners, nil
}
