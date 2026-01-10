package main

import (
	"fmt"
)

func main() {
	_, _ = fmt.Print("Hello World", "\n")
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			s := ""
			switch 0 {
			case (i*3 + j) % 5:
				s = "\033[1B v"
			case (i*3 + j) % 12:
				s = "\033[2B  w"
			}
			_, _ = fmt.Print(s, "texttexttexttexttext text ", i, " ", j, " ", i*3+j, "\n")
		}
	}
	_, _ = fmt.Print("End", "\n")
}
