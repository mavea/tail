package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "alt-save\n")
	_, _ = fmt.Fprint(os.Stdout, "\0337")
	_, _ = fmt.Fprint(os.Stdout, "alt-move\n")
	_, _ = fmt.Fprint(os.Stdout, "\0338")
	_, _ = fmt.Fprint(os.Stdout, "alt-restore\n")
}
