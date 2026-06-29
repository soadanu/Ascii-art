package main

import (
	"fmt"
	"os"
	"strings"
	"ascii-art/functions"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3{
		fmt.Println("Invalid run command")
		fmt.Println("Usage: go run . [string] [font]")
		return
	}

	data := "standard.txt"

	if len(os.Args) == 3 {
		data = os.Args[2]

		// Add .txt only if it isn't already there
		if !strings.HasSuffix(data, ".txt") {
			data += ".txt"
		}
	}

	char, err := functions.LoadBanner(data)
	if err != nil {
		switch err.Error() {
		case "empty banner file":
			fmt.Println("empty banner file")
		case "incomplete banner file":
			fmt.Println("incomplete banner file")
		default:
			fmt.Println("Invalid command")
		}
		return
	}

	fmt.Print(functions.Render(os.Args[1], char))
}
