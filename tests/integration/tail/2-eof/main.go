package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "Hello World", "\n")
	_, _ = fmt.Fprint(os.Stdout, "and", "\n")
	_, _ = fmt.Fprint(os.Stdout, "End")
}
