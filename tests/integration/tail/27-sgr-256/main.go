package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "\033[38;5;202mfg256\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[48;5;25mbg256\033[0m\n")
}
