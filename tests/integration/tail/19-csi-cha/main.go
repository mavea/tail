package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "abcdef\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[4G")
	_, _ = fmt.Fprint(os.Stdout, "H\n")
}
