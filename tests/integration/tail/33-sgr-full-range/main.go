package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "\033[31mred\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[1mbold\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[4munder\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[38;5;255mfg-255\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[48;5;0mbg-0\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[38;2;255;255;255mfg-rgb-max\033[0m\n")
}
