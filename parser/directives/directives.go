package directives

import (
	"github.com/bradleyjkemp/git-owners/file"
	"github.com/pkg/errors"
)

var Parsers = make(map[string]func(string, *file.OwnersFile) error)

func parseFlag(flagName string, o *file.OwnersFile) error {
	if o.Flags == nil {
		o.Flags = make(map[string]bool)
	}
	o.Flags[flagName] = true
	return nil
}

func parseIgnore(value string, o *file.OwnersFile) error {
	pattern, err := parsePattern(value)
	if err != nil {
		return errors.Wrapf(err, "failed to parse ignore directive")
	}

	o.Ignored = append(o.Ignored, pattern)
	return nil
}

func parseWatcher(value string, o *file.OwnersFile) error {
	watcher, err := ParseOwner(value)
	if err != nil {
		return err
	}

	o.Watchers = append(o.Watchers, watcher)
	return nil
}

func init() {
	Parsers["set"] = parseFlag
	Parsers["ignore"] = parseIgnore
	Parsers["watcher"] = parseWatcher
}
