package main

import (
	"strings"
)

// split the input into "words"
func cleanInput(input string) []string {
	return strings.Fields(strings.ToLower(input))
}
