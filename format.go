package env

import (
	"strings"
	"unicode"
)

// formatName formats the name of a struct field to be used as an env variable name
// examples:
// - "DBName" -> "DB_NAME"
// - "DB" -> "DB"
// - "Name" -> "NAME"
// - "AllowedOrigins" -> "ALLOWED_ORIGINS"
// - "TestA" -> "TEST_A"
// - "UserID" -> "USER_ID"
func formatName(name string) string {
	if len(name) == 0 {
		return ""
	}

	var result strings.Builder
	l := len(name) - 1

	result.WriteRune(rune(name[0]))
	for i := 1; i < len(name); i++ {
		r := rune(name[i])
		if unicode.IsUpper(r) && (unicode.IsLower(rune(name[i-1])) ||
			(i < l && unicode.IsLower(rune(name[i+1])))) {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToUpper(result.String())
}
