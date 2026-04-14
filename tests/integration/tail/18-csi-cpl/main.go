package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "q1\n")
	_, _ = fmt.Fprint(os.Stdout, "q2\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[1F")
	_, _ = fmt.Fprint(os.Stdout, "prev\n")
}
