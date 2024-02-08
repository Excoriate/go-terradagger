package utils

import (
	"strings"

	"github.com/google/uuid"
)

func RemoveDoubleQuotes(target string) string {
	return strings.Trim(target, "\"")
}

func GetUUID() string {
	return uuid.New().String()
}

func EscapeValues(value string) string {
	// Perform necessary escaping here. This is a basic example.
	// You might need a more comprehensive approach depending on your input.
	return strings.ReplaceAll(value, "'", "\\'")
}
