package util

import "regexp"

func CompileRegex(expression string) *regexp.Regexp {
	rx, err := regexp.Compile(expression)

	if err != nil {
		panic(err)
	}

	return rx
}
