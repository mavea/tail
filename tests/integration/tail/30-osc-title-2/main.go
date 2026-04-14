package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "start\n")
	_, _ = fmt.Fprint(os.Stdout, "\033]2;osc-two-title\a")
	_, _ = fmt.Fprint(os.Stdout, "end\n")
}
