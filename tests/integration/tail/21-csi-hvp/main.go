package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "A\n")
	_, _ = fmt.Fprint(os.Stdout, "B\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[1;3H")
	_, _ = fmt.Fprint(os.Stdout, "HVP\n")
}
