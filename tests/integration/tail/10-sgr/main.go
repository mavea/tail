package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "normal\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[31mred\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[1;34mbold-blue\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[38;5;202mindexed\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[48;2;1;2;3mtrue-bg\033[0m\n")
}
