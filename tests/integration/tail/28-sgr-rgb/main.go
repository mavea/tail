package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "\033[38;2;10;20;30mfg-rgb\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[48;2;1;2;3mbg-rgb\033[0m\n")
}
