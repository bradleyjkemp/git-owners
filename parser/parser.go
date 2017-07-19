package parser

import (
	"bufio"
	"github.com/bradleyjkemp/git-owners/file"
	"github.com/bradleyjkemp/git-owners/parser/directives"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

var (
	directiveRegex = regexp.MustCompile(`^@([a-z-_]+) (.*)$`)
	groupRegex     = regexp.MustCompile(`^\((.*)\)(?: (.*))?$`)
	commentRegex   = regexp.MustCompile(`^#.*$`)
)

// returns the directive name and remainder if line is a directive, else returns ""
func IsDirective(line string) bool {
	return directiveRegex.Match([]byte(line))
}

func ParseDirective(line string, o *file.OwnersFile) error {
	match := directiveRegex.FindStringSubmatch(line)
	if match == nil {
		return errors.Errorf("failed to parse directive \"%s\"", line)
	}
	directive := match[1]
	value := match[2]

	parser := directives.Parsers[directive]
	if parser == nil {
		return errors.Errorf("unknown directive \"%s\"", directive)
	}

	return parser(value, o)
}

func ParseOwner(line string, o *file.OwnersFile) error {
	owner, err := directives.ParseOwner(line)
	if err != nil {
		return err
	}

	o.Owners = append(o.Owners, owner)
	return nil
}

func shouldSkip(line string) bool {
	if len(line) == 0 {
		return true
	}

	return commentRegex.Match([]byte(line))
}

func ParseFile(input *bufio.Scanner) (*file.OwnersFile, error) {
	var lineCount int
	owners := &file.OwnersFile{}

	for input.Scan() {
		lineCount++
		line := strings.TrimSpace(input.Text())

		if shouldSkip(line) {
			continue
		}

		var err error
		if IsDirective(line) {
			err = ParseDirective(line, owners)
		} else {
			err = ParseOwner(line, owners)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed parsing line %d", lineCount)
		}
	}

	if err := input.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to parse file at line %d", lineCount)
	}

	return owners, nil
}
