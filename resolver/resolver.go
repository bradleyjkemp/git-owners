package resolver

// Given a path to a file this package resolves all the owners of that file

import (
	"bufio"
	"github.com/bradleyjkemp/git-owners/file"
	"github.com/bradleyjkemp/git-owners/git"
	"github.com/bradleyjkemp/git-owners/parser"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func ResolveOwners(path string, allOwners bool) ([]string, error) {
	_, file := filepath.Split(path)

	gitRoot, err := git.RepoRoot()
	if err != nil {
		return nil, err
	}

	return resolveOwners(&resolveOwnersArgs{
		directory:           filepath.Dir(path),
		filename:            file,
		resolveParentOwners: allOwners,
		parser:              parseFile,
		gitRoot:             gitRoot,
	})
}

func ResolveOwnersAtCommit(path string, allOwners bool, commit string) ([]string, error) {
	_, file := filepath.Split(path)

	gitRoot, err := git.RepoRoot()
	if err != nil {
		return nil, err
	}

	return resolveOwners(&resolveOwnersArgs{
		directory:           filepath.Dir(path),
		filename:            file,
		resolveParentOwners: allOwners,
		parser:              parseFileAtCommit(commit),
		gitRoot:             gitRoot,
	})
}

type fileParser func(string) (*file.OwnersFile, error)

type resolveOwnersArgs struct {
	directory           string
	filename            string
	resolveParentOwners bool
	gitRoot             string
	parser              fileParser
	owners              []string
}

func resolveOwners(args *resolveOwnersArgs) ([]string, error) {
	ownersPath := filepath.Join(args.directory, "OWNERS")

	if _, err := os.Stat(ownersPath); os.IsNotExist(err) {
		if args.directory != args.gitRoot {
			args.directory = filepath.Dir(args.directory)
			return resolveOwners(args)
		}
		return args.owners, nil
	}

	ownersFile, err := parseFile(ownersPath)
	if err != nil {
		return nil, err
	}

	// file is ignored so any previous owners are forgotten and we no longer need to look in parent directories
	if isIgnored(args.filename, ownersFile.Ignored) {
		return nil, nil
	}

	if args.resolveParentOwners || len(args.owners) == 0 {
		owners := matchOwners(args.filename, ownersFile.Owners)
		args.owners = append(args.owners, owners...)
	}

	if args.directory == args.gitRoot {
		return args.owners, nil
	}

	// always need to look in the parent directory in case filename is ignored
	// however we may not need to care about keeping track of any further owners
	args.resolveParentOwners = args.resolveParentOwners && !ownersFile.Flags["noparent"]
	args.directory = filepath.Dir(args.directory)
	return resolveOwners(args)
}

func isIgnored(filename string, ignoreRules []string) bool {
	for _, rule := range ignoreRules {
		if matched, _ := filepath.Match(rule, filename); matched {
			return true
		}
	}

	return false
}

func matchOwners(filename string, owners []*file.Owner) []string {
	var matchedOwners []string
	for _, owner := range owners {
		if matched, _ := filepath.Match(owner.Pattern, filename); matched {
			matchedOwners = append(matchedOwners, owner.User)
		}
	}

	return matchedOwners
}

func parseFile(path string) (*file.OwnersFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open OWNERS file %s", path)
	}

	return parser.ParseFile(bufio.NewScanner(file))
}

func parseFileAtCommit(commit string) fileParser {
	return func(path string) (*file.OwnersFile, error) {
		file := git.FileContentsAtCommit(path, commit)
		return parser.ParseFile(bufio.NewScanner(strings.NewReader(file)))
	}
}
