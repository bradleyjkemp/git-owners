package directives

import (
	"github.com/bradleyjkemp/git-owners/file"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
)

func ParseOwner(line string) (*file.Owner, error) {
	tokens := strings.SplitN(line, " ", 3)

	if len(tokens) == 3 {
		return nil, errors.Errorf("failed to parse owner \"%s\", too many components", line)
	}

	if len(tokens) == 1 {
		return &file.Owner{
			User:    tokens[0],
			Pattern: "*",
		}, nil
	} else {
		pattern, err := parsePattern(tokens[1])
		if err != nil {
			return nil, err
		}

		return &file.Owner{
			User:    tokens[0],
			Pattern: pattern,
		}, nil
	}
}

func parsePattern(pattern string) (string, error) {
	// do a dummy match to check the pattern is valid
	_, err := filepath.Match(pattern, "dummy")
	if err == filepath.ErrBadPattern {
		return "", errors.Wrapf(err, "failed to parse pattern \"%s\"", pattern)
	}

	return pattern, nil
}
