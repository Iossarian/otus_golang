package main

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

var b strings.Builder

func main() {
}

func Unpack(incomingString string) (string, error) {
	slicedString := []rune(incomingString)
	stringLen := len(slicedString)
	for i, char := range slicedString {
		prevIndex := i - 1
		nextIndex := i + 1
		if unicode.IsDigit(char) {
			if prevIndex < 0 || unicode.IsDigit(slicedString[nextIndex]) {
				return "", ErrInvalidString
			}
			repeatCount := int(char - '0')
			prevChar := slicedString[i-1]
			b.WriteString(strings.Repeat(string(prevChar), repeatCount))
		} else {
			isLastChar := nextIndex == stringLen
			if nextIndex <= stringLen && (isLastChar || !unicode.IsDigit(slicedString[nextIndex])) {
				b.WriteRune(char)
			}
		}
	}
	result := b.String()
	b.Reset()
	return result, nil
}
