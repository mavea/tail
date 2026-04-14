package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "x1\n")
	_, _ = fmt.Fprint(os.Stdout, "x2\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[2E")
	_, _ = fmt.Fprint(os.Stdout, "next\n")
}
