package commands

import (
	"strings"
)

func cleanInput(text string) []string {
	return strings.Split(strings.Trim(strings.ToLower(text), " "), " ")
}
