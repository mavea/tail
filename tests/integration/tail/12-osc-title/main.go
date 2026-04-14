package main

import (
	"fmt"
	"os"
)

func main() {
	_, _ = fmt.Fprint(os.Stdout, "start\n")
	_, _ = fmt.Fprint(os.Stdout, "\033]0;tail-test-title\a")
	_, _ = fmt.Fprint(os.Stdout, "body\n")
	_, _ = fmt.Fprint(os.Stdout, "\033]2;window-name\a")
	_, _ = fmt.Fprint(os.Stdout, "done\n")
}
