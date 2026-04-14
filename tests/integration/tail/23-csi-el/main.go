package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "line-to-clear\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[2K")
	_, _ = fmt.Fprint(os.Stdout, "line-new\n")
}
