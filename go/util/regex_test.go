package util

import (
	"regexp"
	"testing"
)

func TestValidRegex(t *testing.T) {
	expression := "^\\d+$"

	var rx *regexp.Regexp

	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Fatalf("regex '%s' should compile, got %v instead", expression, err)
			}
		}()

		rx = CompileRegex(expression)
	}()

	if !rx.MatchString("123") {
		t.Errorf("'123' should match '%s'", expression)
	}

	if rx.MatchString("foo") {
		t.Errorf("'foo' should not match '%s'", expression)
	}
}

func TestInvalidRegex(t *testing.T) {
	expression := "("

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("regex '%s' should not compile", expression)
			}
		}()

		CompileRegex(expression)
	}()
}
