package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "r1\n")
	_, _ = fmt.Fprint(os.Stdout, "r2\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[2d")
	_, _ = fmt.Fprint(os.Stdout, "row2\n")
}
