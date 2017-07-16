package resolver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var resolverAllOwners = []struct {
	filepath       string
	expected       [][]string
	expectedDirect [][]string
}{
	{"../example_repo/go/src/test/example.go", [][]string{{"bob", "charlie"}, {"alice"}}, [][]string{{"bob", "charlie"}}},
	{"../example_repo/go/src/test/README", [][]string{{"alice"}}, [][]string{{"alice"}}},
	{"../example_repo/go/BUILD", nil, nil},
	{"../example_repo/BUILD", nil, nil},
	{"../example_repo/js/BUILD", nil, nil},
	{"../example_repo/js/static/test.png", [][]string{{"dave@example.com"}}, [][]string{{"dave@example.com"}}},
}

func TestResolverAllOwners(t *testing.T) {
	for _, testCase := range resolverAllOwners {
		actual, err := ResolveOwners(testCase.filepath, true)
		assert.NoError(t, err, "must not fail to resolve")
		assert.Equal(t, testCase.expected, actual)
	}
}

func TestResolverDirectOwners(t *testing.T) {
	for _, testCase := range resolverAllOwners {
		actual, err := ResolveOwners(testCase.filepath, false)
		assert.NoError(t, err, "must not fail to resolve")
		assert.Equal(t, testCase.expectedDirect, actual)
	}
}
