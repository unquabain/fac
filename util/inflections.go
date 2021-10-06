package util

import (
	"strings"
)

func Parameterize(input string) string {
	return strings.ToLower(
		strings.ReplaceAll(input, ` `, `-`),
	)
}
