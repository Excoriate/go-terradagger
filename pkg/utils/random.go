package utils

import (
	"math/rand"
	"strings"
	"unicode/utf8"

	"github.com/docker/docker/pkg/namesgenerator"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytesWithSpecialChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()"

type RandomStringOptions struct {
	Length            int
	AllowSpecialChars bool
	UpperLowerRandom  bool
}

func GenerateRandomString(opt RandomStringOptions) string {
	var letterBytesRunes []rune
	if opt.AllowSpecialChars {
		letterBytesRunes = []rune(letterBytesWithSpecialChars)
	} else {
		letterBytesRunes = []rune(letterBytes)
	}

	b := make([]rune, opt.Length)
	for i := range b {
		b[i] = letterBytesRunes[rand.Intn(len(letterBytesRunes))]
	}

	if opt.UpperLowerRandom {
		for i, v := range b {
			if rand.Intn(2) == 0 {
				upperRune, _ := utf8.DecodeRuneInString(strings.ToUpper(string(v)))
				b[i] = upperRune
			} else {
				lowerRune, _ := utf8.DecodeRuneInString(strings.ToLower(string(v)))
				b[i] = lowerRune
			}
		}
	}

	return string(b)
}

func GenerateRandomName(retry int) string {
	return namesgenerator.GetRandomName(retry)
}
