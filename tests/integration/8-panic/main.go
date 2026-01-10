package main

import (
	"fmt"
)

func main() {
	var sl []string
	_, _ = fmt.Print("Hello World", "\n")
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			_, _ = fmt.Print("texttexttexttexttext text ", i, " ", j, " ", i*3+j, "\n")
		}
	}
	_, _ = fmt.Print("Error", "\n")

	// #nosec G602
	sl[1] = ""
}
