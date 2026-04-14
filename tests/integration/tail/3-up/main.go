package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "Hello World", "\n")
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			s := ""
			switch 0 {
			case (i*3 + j) % 5:
				s = "\033[1A v"
			case (i*3 + j) % 12:
				s = "\033[2A  w"
			}
			_, _ = fmt.Fprint(os.Stdout, s, "texttexttexttexttext text ", i, " ", j, " ", i*3+j, "\n")
		}
	}
	_, _ = fmt.Fprint(os.Stdout, "End", "\n")
}
