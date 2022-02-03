package enthstore

import (
	"strings"
)

func quoteKey(key string) string {
	return `'` + strings.ReplaceAll(key, "'", `''`) + `'`
}

func quoteValue(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	return `"` + value + `"`
}
