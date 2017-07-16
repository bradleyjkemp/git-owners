package parser

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	// "github.com/davecgh/go-spew/spew"
)

func flagSet(flags ...string) map[string]bool {
	set := make(map[string]bool)
	for _, flag := range flags {
		set[flag] = true
	}
	return set
}

func owners(users ...string) *Owner {
	return &Owner{
		Users:   users,
		Pattern: "*",
	}
}

func patternOwners(pattern string, users ...string) *Owner {
	return &Owner{
		Users:   users,
		Pattern: pattern,
	}
}

var positiveParserTests = []struct {
	file     string
	expected *OwnersFile
}{
	{
		`@set noparent
		alice
		(bob & carol) *.js`,
		&OwnersFile{
			Flags: flagSet("noparent"),
			Owners: []*Owner{
				owners("alice"),
				patternOwners("*.js", "bob", "carol"),
			},
		},
	},
	{
		`
		# comments and blank lines are ignored
		
		alice123@example.com
		(single user group is okay...)
		`,
		&OwnersFile{
			Owners: []*Owner{
				owners("alice123@example.com"),
				owners("single user group is okay..."),
			},
		},
	},
	{
		`@set noparent
		alice
		(bob & carol) *.js`,
		&OwnersFile{
			Flags: flagSet("noparent"),
			Owners: []*Owner{
				owners("alice"),
				patternOwners("*.js", "bob", "carol"),
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
