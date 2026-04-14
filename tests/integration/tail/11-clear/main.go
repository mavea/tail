package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "line-1\n")
	_, _ = fmt.Fprint(os.Stdout, "line-2\n")
	_, _ = fmt.Fprint(os.Stdout, "line-3\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[1A\033[2K")
	_, _ = fmt.Fprint(os.Stdout, "line-2-updated\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[2J")
	_, _ = fmt.Fprint(os.Stdout, "after-clear\n")
}
