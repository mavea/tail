package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "aa\n")
	_, _ = fmt.Fprint(os.Stdout, "bb\n")
	_, _ = fmt.Fprint(os.Stdout, "cc\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[2A")
	_, _ = fmt.Fprint(os.Stdout, "UP\n")
}
