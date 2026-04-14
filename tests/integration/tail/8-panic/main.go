package main

import (
	"fmt"
	"os"
)

func main() {
	var sl []string
	_, _ = fmt.Fprint(os.Stdout, "Hello World", "\n")
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			_, _ = fmt.Fprint(os.Stdout, "texttexttexttexttext text ", i, " ", j, " ", i*3+j, "\n")
		}
	}
	_, _ = fmt.Fprint(os.Stdout, "Error", "\n")

	// #nosec G602
	sl[1] = ""
}
