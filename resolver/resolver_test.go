package resolver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var resolveOwnersTests = []struct {
	filepath       string
	expected       []string
	expectedDirect []string
}{
	{"../example_repo/go/src/test/example.go",
		[]string{"bob", "carol", "alice", "bradleyjkemp"}, []string{"bob", "carol"}},

	// should only own files where the pattern matches
	{"../example_repo/go/src/test/README",
		[]string{"carol", "alice", "bradleyjkemp"}, []string{"carol"}},

	// ignored files should always be ignored even if there is an owner lower down
	{"../example_repo/go/BUILD", nil, nil},
	{"../example_repo/BUILD", nil, nil},
	{"../example_repo/js/BUILD", nil, nil},

	// should respect the @set noparent flag
	{"../example_repo/js/static/test.png",
		[]string{"dave@example.com"}, []string{"dave@example.com"}},
}

func TestResolverAllOwners(t *testing.T) {
	for _, testCase := range resolveOwnersTests {
		actual, err := ResolveOwners(testCase.filepath, true)
		assert.NoError(t, err, "must not fail to resolve")
		assert.Equal(t, testCase.expected, actual, testCase.filepath)
	}
}

func TestResolverDirectOwners(t *testing.T) {
	for _, testCase := range resolveOwnersTests {
		actual, err := ResolveOwners(testCase.filepath, false)
		assert.NoError(t, err, "must not fail to resolve")
		assert.Equal(t, testCase.expectedDirect, actual, testCase.filepath)
	}
}
