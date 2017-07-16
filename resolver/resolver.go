package resolver

// Given a path to a file this package resolves all the owners of that file

import (
	"bufio"
	"github.com/bradleyjkemp/git-owners/git"
	"github.com/bradleyjkemp/git-owners/parser"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func ResolveOwners(path string, allOwners bool) ([][]string, error) {
	_, file := filepath.Split(path)
	return resolveOwners(filepath.Dir(path), file, allOwners, parseFile)
}

func ResolveOwnersAtCommit(path string, allOwners bool, commit string) ([][]string, error) {
	_, file := filepath.Split(path)
	return resolveOwners(filepath.Dir(path), file, allOwners, parseFileAtCommit(commit))
}

type fileParser func(string) (*parser.OwnersFile, error)

func resolveOwners(directory, filename string, resolveParents bool, parser fileParser) ([][]string, error) {
	gitRoot, err := git.RepoRoot()
	if err != nil {
		return nil, err
	}

	ownersPath := filepath.Join(directory, "OWNERS")

	if _, err = os.Stat(ownersPath); os.IsNotExist(err) {
		if directory != gitRoot {
			return resolveOwners(filepath.Dir(directory), filename, resolveParents, parser)
		}
		return nil, nil
	}

	ownersFile, err := parseFile(ownersPath)
	if err != nil {
		return nil, err
	}

	if isIgnored(filename, ownersFile.Ignored) {
		return nil, nil
	}

	owners := matchOwners(filename, ownersFile.Owners)

	if directory == gitRoot || ownersFile.Flags["noparent"] {
		return owners, nil
	}

	if resolveParents || len(owners) == 0 {
		parentOwners, err := resolveOwners(filepath.Dir(directory), filename, resolveParents, parser)
		if err != nil {
			return nil, err
		}

		return append(owners, parentOwners...), nil
	}

	return owners, nil
}

func isIgnored(filename string, ignoreRules []string) bool {
	for _, rule := range ignoreRules {
		if matched, _ := filepath.Match(rule, filename); matched {
			return true
		}
	}

	return false
}

func matchOwners(filename string, owners []*parser.Owner) [][]string {
	var matchedOwners [][]string
	for _, owner := range owners {
		if matched, _ := filepath.Match(owner.Pattern, filename); matched {
			matchedOwners = append(matchedOwners, owner.Users)
		}
	}

	return matchedOwners
}

func parseFile(path string) (*parser.OwnersFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open OWNERS file %s", path)
	}

	return parser.ParseFile(bufio.NewScanner(file))
}

func parseFileAtCommit(commit string) fileParser {
	return func(path string) (*parser.OwnersFile, error) {
		file := git.FileContentsAtCommit(path, commit)
		return parser.ParseFile(bufio.NewScanner(strings.NewReader(file)))
	}
}
