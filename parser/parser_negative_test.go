package parser

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var negativeParserTests = []struct {
	file          string
	expectedError string
}{
	{
		`@noparent true
		alice`,
		"failed parsing line 1: unknown directive \"noparent\"",
	},
	{
		`@set noparent
		bob *.go *.js`,
		"failed parsing line 2: failed to parse owner \"bob *.go *.js\", too many components",
	},
	{
		`carol [a-]`,
		"failed parsing line 1: failed to parse pattern \"[a-]\": syntax error in pattern",
	},
	{
		`@watcher dave [b-]
		@watcher eve`,
		"failed parsing line 1: failed to parse pattern \"[b-]\": syntax error in pattern",
	},
}

func TestParserNegative(t *testing.T) {
	for _, testCase := range negativeParserTests {
		scanner := bufio.NewScanner(strings.NewReader(testCase.file))
		_, err := ParseFile(scanner)
		assert.EqualError(t, err, testCase.expectedError, testCase.file)
	}
}
