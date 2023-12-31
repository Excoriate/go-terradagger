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

func CleanSliceFromValuesThatAreEmpty(slice []string) []string {
	var cleaned []string
	for _, value := range slice {
		if value != "" {
			cleaned = append(cleaned, value)
		}
	}
	return cleaned
}

func MisSlices(slices ...[]string) []string {
	var mixed []string
	for _, slice := range slices {
		if slice == nil {
			continue
		}
		mixed = append(mixed, slice...)
	}
	return mixed
}

func MixSlices(slices ...[]string) []string {
	var mixed []string
	for _, slice := range slices {
		mixed = append(mixed, slice...)
	}
	return mixed
}

func MergeMaps(maps ...map[string]string) map[string]string {
	mixed := make(map[string]string)
	for _, maps := range maps {
		for key, value := range maps {
			// normalise the key and the value by removing double quotes
			mixed[RemoveDoubleQuotes(key)] = RemoveDoubleQuotes(value)
		}
	}
	return mixed
}
