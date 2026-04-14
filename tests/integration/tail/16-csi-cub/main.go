package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "ABCDE\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[2D")
	_, _ = fmt.Fprint(os.Stdout, "L\n")
}
