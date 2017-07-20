package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var resolveOwners = []struct {
	filepath       string
	expected       string
	expectedDirect string
}{
	{"example_repo/go/src/test/example.go",
		"example_repo/go/src/test/example.go : [bob carol alice bradleyjkemp]\n",
		"example_repo/go/src/test/example.go : [bob carol]\n",
	},

	// should only own files where the pattern matches
	{"example_repo/go/src/test/README",
		"example_repo/go/src/test/README : [carol alice bradleyjkemp]\n",
		"example_repo/go/src/test/README : [carol]\n",
	},

	// ignored files should always be ignored even if there is an owner lower down
	{"example_repo/go/BUILD",
		"example_repo/go/BUILD : []\n",
		"example_repo/go/BUILD : []\n",
	},
	{"example_repo/js/BUILD",
		"example_repo/js/BUILD : []\n",
		"example_repo/js/BUILD : []\n",
	},
	{"example_repo/BUILD",
		"example_repo/BUILD : []\n",
		"example_repo/BUILD : []\n",
	},

	// should respect the @set noparent flag
	{"example_repo/js/static/test.png",
		"example_repo/js/static/test.png : [dave@example.com]\n",
		"example_repo/js/static/test.png : [dave@example.com]\n",
	},
}

func TestResolverAllOwners(t *testing.T) {
	var buf bytes.Buffer
	stdout = &buf

	for _, testCase := range resolveOwners {
		buf.Reset()
		fileOwners(cliFlags{
			"master",
			true,
			[]string{testCase.filepath},
		})
		assert.Equal(t, testCase.expected, buf.String(), testCase.filepath)
	}
}

func TestResolverDirectOwners(t *testing.T) {
	var buf bytes.Buffer
	stdout = &buf

	for _, testCase := range resolveOwners {
		buf.Reset()
		fileOwners(cliFlags{
			"master",
			false,
			[]string{testCase.filepath},
		})
		assert.Equal(t, testCase.expectedDirect, buf.String(), testCase.filepath)
	}
}
