package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "Hello World", "\n")
	_, _ = fmt.Fprint(os.Stdout, "12345", "\n")
	_, _ = fmt.Fprint(os.Stdout, "67890", "\n")
	_, _ = fmt.Fprint(os.Stdout, "abcde", "\n")

	_, _ = fmt.Fprint(os.Stdout, "*\033[2Ag\033[2Ch\033[1Bi\033[3Dj")
	_, _ = fmt.Fprint(os.Stdout, "\033[2Fk\033[s\033[1El\033[2C\033[um")

	_, _ = fmt.Fprint(os.Stdout, "\n", "\n", "\n")
	_, _ = fmt.Fprint(os.Stdout, "\033[1;0Hfg\033[1;1H-\033[1;3fe")
}
