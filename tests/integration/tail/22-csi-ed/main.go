package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "before\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[2J")
	_, _ = fmt.Fprint(os.Stdout, "after\n")
}
