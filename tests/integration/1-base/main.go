package main

import (
	"fmt"
	"time"
)

func main() {
	_, _ = fmt.Print("Hello World", "\n")
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			_, _ = fmt.Print("texttexttexttexttext text ", i, " ", j, " ", i*3+j, "\n")
		}
		time.Sleep(time.Second / 3)
	}
	time.Sleep(time.Second / 2)
	_, _ = fmt.Print("End", "\n")
}
