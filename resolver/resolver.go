package resolver

// Given a path to a file this package resolves all the owners of that file

import (
	"bufio"
	"github.com/bradleyjkemp/git-owners/git"
	"github.com/bradleyjkemp/git-owners/parser"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// Resolves all owners up to the root of the repo
func ResolveAllOwners(path string) ([][]string, error) {
    directory, file := filepath.Split(path)
	return resolveOwners(directory, file, true)
}

// Only resolves the most specific owners, stops after the first OWNERS file encountered
func ResolveDirectOwners(path string) ([][]string, error) {
    directory, file := filepath.Split(path)
	return resolveOwners(directory, file, false)
}

func resolveOwners(directory, filename string, resolveParents bool) ([][]string, error) {
	gitRoot, err := git.RepoRoot()
	if err != nil {
		return nil, err
	}

	ownersPath := filepath.Join(directory, "OWNERS")

    if _, err = os.Stat(ownersPath); os.IsNotExist(err) {
        if directory != gitRoot {
            return resolveOwners(filepath.Dir(directory), filename, resolveParents)
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

	if resolveParents && directory == gitRoot && !ownersFile.Flags["noparent"] {
        parentOwners, err := resolveOwners(filepath.Dir(directory), filename, resolveParents)
        if err != nil {
            return nil, err
        }
        
        return append(owners, parentOwners), nil
	}
	
	return owners, nil
}

func isIgnored(filename string, ignoreRules []string) bool {
    for _, rule := ignoreRules {
        if filePath.Match(filename, rule) {
            return true
        }
    }
    
    return false
}

func matchOwners(filename string, owners []*parser.Owners) [][]string {
    var matchedOwners [][]string
    for _, owner := range owners {
        if filepath.Match(filename, owners.Pattern) {
            matchedOwners = append(matchedOwners, owners.Users)
        }
    }
}

func parseFile(path string) (*parser.OwnersFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open OWNERS file %s", path)
	}

	return parser.ParseFile(bufio.NewScanner(file))
}
