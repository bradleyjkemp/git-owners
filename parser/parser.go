package parser

import (
	"bufio"
	"github.com/pkg/errors"
	"path/filepath"
	"regexp"
	"strings"
)

type Owner struct {
	Users   []string
	Pattern string
}

type OwnersFile struct {
	Flags    map[string]bool
	Ignored  []string
	Owners   []*Owner
	Watchers []*Owner
}

var (
	directiveRegex = regexp.MustCompile(`^@([a-z-_]+) (.*)$`)
	groupRegex     = regexp.MustCompile(`^\((.*)\)(?: (.*))?$`)
	commentRegex   = regexp.MustCompile(`^#.*$`)
)

// returns the directive name and remainder if line is a directive, else returns ""
func IsDirective(line string) bool {
	return directiveRegex.Match([]byte(line))
}

func (o *OwnersFile) ParseDirective(line string) error {
	match := directiveRegex.FindStringSubmatch(line)
	if match == nil {
		return errors.Errorf("failed to parse directive \"%s\"", line)
	}
	directive := match[1]
	value := match[2]

	switch directive {
	case "set":
		if o.Flags == nil {
			o.Flags = make(map[string]bool)
		}
		o.Flags[value] = true

	case "ignore":
		pattern, err := parsePattern(value)
		if err != nil {
			return errors.Wrapf(err, "failed to parse ignore directive")
		}

		o.Ignored = append(o.Ignored, pattern)

	case "watchers":

	default:
		return errors.Errorf("unknown directive \"%s\"", directive)
	}

	return nil
}

func (o *OwnersFile) ParseOwner(line string) error {
	owner, err := parseOwner(line)
	if err != nil {
		return err
	}

	o.Owners = append(o.Owners, owner)
	return nil
}

func parseOwner(line string) (*Owner, error) {
	if groupRegex.Match([]byte(line)) {
		return parseGroupOwners(line)
	}

	tokens := strings.SplitN(line, " ", 2)

	if len(tokens) == 1 {
		return &Owner{
			Users:   []string{tokens[0]},
			Pattern: "*",
		}, nil
	} else {
		pattern, err := parsePattern(tokens[1])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse owner \"%s\"", line)
		}

		return &Owner{
			Users:   []string{tokens[0]},
			Pattern: pattern,
		}, nil
	}
}

func parseGroupOwners(line string) (*Owner, error) {
	match := groupRegex.FindStringSubmatch(line)
	if match == nil {
		return nil, errors.Errorf("error parsing group \"%s\"", line)
	}

	groupString := match[1]
	owners := strings.Split(groupString, " & ")

	var pattern string
	if match[2] == "" { // if there was no pattern an empty match is still returned
		pattern = "*"
	} else {
		var err error
		pattern, err = parsePattern(match[2])
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing group \"%s\"", groupString)
		}
	}

	return &Owner{
		Users:   owners,
		Pattern: pattern,
	}, nil
}

func parsePattern(input string) (string, error) {
	// do a dummy match to check the pattern is valid
	_, err := filepath.Match(input, "")
	if err == filepath.ErrBadPattern {
		return "", errors.Wrapf(err, "failed to parse pattern %s", input)
	}

	return input, nil
}

func shouldSkip(line string) bool {
	if len(line) == 0 {
		return true
	}

	return commentRegex.Match([]byte(line))
}

func ParseFile(file *bufio.Scanner) (*OwnersFile, error) {
	var lineCount int
	owners := &OwnersFile{}

	for file.Scan() {
		lineCount++
		line := strings.TrimSpace(file.Text())

		if shouldSkip(line) {
			continue
		}

		var err error
		if IsDirective(line) {
			err = owners.ParseDirective(line)
		} else {
			err = owners.ParseOwner(line)
		}

		if err != nil {
			return nil, errors.Wrapf(err, "failed parsing line %d", lineCount)
		}
	}

	if err := file.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to parse file at line %d", lineCount)
	}

	return owners, nil
}
