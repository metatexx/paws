package main

import "regexp"

func highlight(input string, re *regexp.Regexp) string {
	indices := re.FindStringSubmatchIndex(input)
	if indices == nil {
		return input
	}

	highlighted := input[:indices[0]] +
		"\033[1;31m" + // Start red highlight
		input[indices[0]:indices[1]] +
		"\033[0m" + // End highlight
		input[indices[1]:]

	return highlighted
}
