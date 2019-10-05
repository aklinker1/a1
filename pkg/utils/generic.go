package utils

import (
	"fmt"
	"net"
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

// GetOutboundIP -
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
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

// LogWhite -
func LogWhite(template string, args ...interface{}) {
	if IsVerbose() {
		fmt.Printf("    %s\n", fmt.Sprintf(template, args...))
	}
}

// LogRed -
func LogRed(template string, args ...interface{}) {
	if IsVerbose() {
		fmt.Printf("    \x1b[91m%s\x1b[0m\n", fmt.Sprintf(template, args...))
	}
}
