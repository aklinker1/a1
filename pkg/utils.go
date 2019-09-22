package pkg

import (
	"fmt"
	"os"
	"unicode"
)

func pluralize(count int, singular string, plural string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}
	return fmt.Sprintf("%d %s", count, plural)
}

func lowerFirstChar(str string) string {
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
		fmt.Printf("  \x1b[2m%s\x1b[0m\n", fmt.Sprintf(template, args...))
	}
}
