package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "Hello World", "\n")
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			_, _ = fmt.Fprint(os.Stdout, "texttexttexttexttext text ", i, " ", j, " ", i*3+j, "\n")
		}
		time.Sleep(time.Second / 3)
	}
	time.Sleep(time.Second / 2)
	_, _ = fmt.Fprint(os.Stdout, "End", "\n")
}
