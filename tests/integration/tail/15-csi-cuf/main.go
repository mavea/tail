package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "12345\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[3C")
	_, _ = fmt.Fprint(os.Stdout, "R\n")
}
