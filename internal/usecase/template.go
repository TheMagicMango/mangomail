// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package usecase

import (
	"fmt"
	"regexp"
	"strings"
)

var placeholderRegex = regexp.MustCompile(`\{\{(\s*\w+\s*)\}\}`)

// replacePlaceholders replaces template placeholders with actual values from the data map.
// Placeholders use the syntax {{key}} where key matches a key in the data map.
// Whitespace around the key is ignored, so {{name}} and {{ name }} are equivalent.
func replacePlaceholders(text string, data map[string]interface{}) string {
	return placeholderRegex.ReplaceAllStringFunc(text, func(match string) string {

		fieldName := strings.TrimSpace(match[2 : len(match)-2])

		if value, ok := data[fieldName]; ok {
			return fmt.Sprintf("%v", value)
		}

		return match
	})
}
