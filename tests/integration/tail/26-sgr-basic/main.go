package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "\033[31mred\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[1mbold\033[0m\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[4munder\033[0m\n")
}
