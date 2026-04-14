package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "save\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[s")
	_, _ = fmt.Fprint(os.Stdout, "move\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[u")
	_, _ = fmt.Fprint(os.Stdout, "restore\n")
}
