package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "\033[32mgreen\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "12345\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[1A\033[2K")
	_, _ = fmt.Fprint(os.Stdout, "rewritten\n")
}
