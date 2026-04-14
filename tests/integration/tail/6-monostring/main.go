package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "Hello World")
	_, _ = fmt.Fprint(os.Stdout, "12345")
	_, _ = fmt.Fprint(os.Stdout, "67890")
	_, _ = fmt.Fprint(os.Stdout, "abcde")
	_, _ = fmt.Fprint(os.Stdout, "fghijk")
	_, _ = fmt.Fprint(os.Stdout, "lmnopq")
	_, _ = fmt.Fprint(os.Stdout, "rstuvwx")
	_, _ = fmt.Fprint(os.Stdout, "Buy World")
}
