package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "start\n")
	_, _ = fmt.Fprint(os.Stdout, "\033]0;osc-zero-title\a")
	_, _ = fmt.Fprint(os.Stdout, "end\n")
}
