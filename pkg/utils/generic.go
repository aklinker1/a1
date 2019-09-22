package utils

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Pluralize -
func Pluralize(count int, singular string, plural string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}
	return fmt.Sprintf("%d %s", count, plural)
}

// AddS -
func AddS(str string) string {
	if strings.HasSuffix(str, "s") {
		return str
	}
	return str + "s"
}

// LowerFirstChar -
func LowerFirstChar(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// IsVerbose - Whether or not to print output
func IsVerbose() bool {
	return os.Getenv("VERBOSE") == "true"
}

// Log -
func Log(template string, args ...interface{}) {
	if IsVerbose() {
		fmt.Printf("    \x1b[2m%s\x1b[0m\n", fmt.Sprintf(template, args...))
	}
}
