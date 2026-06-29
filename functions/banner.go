package functions

import (
	"errors"
	"os"
	"strings"
)

func LoadBanner(filename string) (map[rune][]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("empty banner file")
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")

	if len(lines) < 855 {
		return nil, errors.New("incomplete banner file")
	}

	chars := make(map[rune][]string)

	for i := 0; i < 95; i++ {
		ch := rune(32 + i)

		start := 1 + i*9
		end := start + 8

		if end > len(lines) {
			return nil, errors.New("incomplete banner file")
		}

		chars[ch] = lines[start:end]
	}

	return chars, nil
}
