package main

import (
	"strings"
	"unicode"
)

func ParseCommand(command string) []string {
	var res []string
	var token strings.Builder
	var insideQuotes bool
	insideToken := false

	command = strings.TrimSpace(command)

	for _, c := range command {
		switch {
		case unicode.IsSpace(c):
			if insideQuotes {
				token.WriteRune(c)
				continue
			}

			if insideToken {
				insideToken = false
				res = append(res, token.String())
				token.Reset()
			}
		case c == '"':
			insideQuotes = !insideQuotes
			if !insideQuotes {
				res = append(res, token.String())
				token.Reset()
				insideToken = false
			}
		default:
			insideToken = true
			token.WriteRune(c)
		}
	}

	if token.Len() > 0 {
		res = append(res, token.String())
		token.Reset()
	}

	return res
}
