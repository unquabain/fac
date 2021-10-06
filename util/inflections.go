package util

import (
	"strings"
)

// Parameterize is like the Ruby ActiveSupport method of the
// same name (though maybe not as smart). It boils down a
// plain English string into an identifier.
func Parameterize(input string) string {
	return strings.ToLower(
		strings.ReplaceAll(input, ` `, `-`),
	)
}
