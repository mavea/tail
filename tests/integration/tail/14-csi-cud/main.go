package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "top\n")
	_, _ = fmt.Fprint(os.Stdout, "mid\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[1B")
	_, _ = fmt.Fprint(os.Stdout, "down\n")
}
