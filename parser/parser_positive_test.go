package parser

import (
	"bufio"
	"github.com/bradleyjkemp/git-owners/file"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func flagSet(flags ...string) map[string]bool {
	set := make(map[string]bool)
	for _, flag := range flags {
		set[flag] = true
	}
	return set
}

var positiveParserTests = []struct {
	file     string
	expected *file.OwnersFile
}{
	{
		`@set noparent
		alice
		bob *.js
		carol ./*.js`,
		&file.OwnersFile{
			Flags: flagSet("noparent"),
			Owners: []*file.Owner{
				{"alice", "*"},
				{"bob", "*.js"},
				{"carol", "./*.js"},
			},
		},
	},
	{
		`
		# comments and blank lines are ignored
		
		alice123@example.com
		`,
		&file.OwnersFile{
			Owners: []*file.Owner{
				{"alice123@example.com", "*"},
			},
		},
	},
	{
		`
		@watcher dave *.go
		`,
		&file.OwnersFile{
			Watchers: []*file.Owner{
				{"dave", "*.go"},
			},
		},
	},
}

func TestParserPositive(t *testing.T) {
	for _, testCase := range positiveParserTests {
		scanner := bufio.NewScanner(strings.NewReader(testCase.file))
		actual, err := ParseFile(scanner)
		assert.NoError(t, err, "should not fail to parse positive test case")
		assert.Equal(t, testCase.expected, actual)
	}
}
