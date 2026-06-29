package functions

import (
	"strings"
	"fmt"
)

func Render(input string, chars map[rune][]string) string {
	if input == "" {
		return ""
	}
	if input == "\\n" {
		fmt.Println()
		return ""
	}

	parts := strings.Split(input, "\\n")

	var result strings.Builder

	for _, part := range parts {
		if part == "" {
			result.WriteString("\n")
			continue
		}

		for row := 0; row < 8; row++ {
			for _, ch := range part {
				if art, ok := chars[ch]; ok {
					result.WriteString(art[row])
				}
			}
			result.WriteString("\n")
		}
	}

	return result.String()
}
