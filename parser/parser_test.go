package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var detectDirectiveTests = []struct {
	line      string
	directive string
	remainder string
}{
	{"@set flagname", "set", "flagname"},
	{"@set this is a flag", "set", "this is a flag"},
	{"@ignore *_test.go", "ignore", "*_test.go"},
	{"@watchers (alice & bob & carol)", "watchers", "(alice & bob & carol)"},
	{"@watchers (alice & bob) *.go", "watchers", "(alice & bob) *.go"},
	{"alice *.go", "", "alice *.go"},
	{"(alice & bob) *.js", "", "(alice & bob) *.js"},
}

func TestDirectives(t *testing.T) {
	for _, testCase := range detectDirectiveTests {
		actualDirective, actualRemainder := IsDirective(testCase.line)

		assert.Equal(t, testCase.directive, actualDirective, testCase.line)
		assert.Equal(t, testCase.remainder, actualRemainder, testCase.line)
	}
}
