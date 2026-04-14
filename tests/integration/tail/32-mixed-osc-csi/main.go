package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "head\n")
	_, _ = fmt.Fprint(os.Stdout, "\033]0;mix-title\a")
	_, _ = fmt.Fprint(os.Stdout, "\033[2J")
	_, _ = fmt.Fprint(os.Stdout, "tail\n")
}
