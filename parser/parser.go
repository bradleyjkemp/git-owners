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
	Owners   []Owner
	Watchers []Owner
}

var (
	directiveRegex = regexp.MustCompile("^@([a-z]+) (.*)$")
	groupRegex     = regexp.MustCompile("^\\((.*)\\) (.*)$")
)

// returns the directive name and remainder if line is a directive, else returns ""
func IsDirective(line string) (directive string, remainder string) {
	match := directiveRegex.FindStringSubmatch(line)
	if match == nil || len(match) != 3 {
		return "", line
	}

	return match[1], match[2]
}

func ParseOwner(line string) (*Owner, error) {
	if groupRegex.Match(line) {
		return ParseGroupOwners(line)
	}

	tokens := strings.SplitN(line, " ", 2)

	if len(tokens) == 1 {
		return &Owner{
			Users:   []string{tokens[0]},
			Pattern: "*",
		}
	} else {
		pattern, err := ParsePattern(tokens[1])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse owner \"%s\"", line)
		}

		return &Owner{
			Users:   []string{tokens[0]},
			Pattern: pattern,
		}
	}
}

func ParseGroupOwners(line string) (*Owner, error) {
	match := groupRegex.FindStringSubmatch(line)
	if match == nil || len(match) != 3 {
		return nil, errors.Errorf("error parsing group \"%s\"", line)
	}

	groupString := match[1]
	owners := strings.Split(groupString, " & ")

	pattern, err := ParsePattern(match[2])
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing group \"%s\"", groupString)
	}

	return Owner{
		Users:   owners,
		Pattern: pattern,
	}, nil
}

func ParsePattern(input string) (string, error) {
	// do a dummy match to check the pattern is valid
	_, err := filepath.Match(value, "")
	if err == filepath.ErrBadPattern {
		return "", errors.Wrapf(err, "failed to parse pattern %s", value)
	}

	return input, nil
}

func ParseDirective(directive string, value string, owners *OwnersFile) error {
	switch directive {
	case "set":
		owners.Flags[value] = true

	case "ignore":
		pattern, err := ParsePattern(value)
		if err != nil {
			return errors.Wrapf(err, "failed to parse ignore directive")
		}

		owners.Ignored = append(owners.Ignored, value)

	case "watchers":

	default:
		return errors.Errorf("unknown directive \"%s\"", directive)
	}

	return nil
}

func ParseFile(file bufio.Scanner) (*OwnersFile, error) {
	var lineCount int
	//owners := &OwnersFile{}

	for file.Scan() {
		lineCount++
		file.Text()
	}

	if err := file.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to parse file at line %d", lineCount)
	}

	return nil, nil
}
