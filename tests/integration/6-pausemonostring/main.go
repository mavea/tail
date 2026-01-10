package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	w := bufio.NewWriter(os.Stdout)
	_, _ = fmt.Fprint(w, "Hello World")
	_ = w.Flush()
	_, _ = fmt.Fprint(w, "12345")
	_ = w.Flush()
	_, _ = fmt.Fprint(w, "67890")
	_ = w.Flush()
	_, _ = fmt.Fprint(w, "abcde")
	_ = w.Flush()
	_, _ = fmt.Fprint(w, "fghijk")
	_ = w.Flush()
	time.Sleep(time.Second / 2)
	_, _ = fmt.Fprint(w, "lmnopq")
	_ = w.Flush()
	_, _ = fmt.Fprint(w, "rstuvwx")
	_ = w.Flush()
	_, _ = fmt.Fprint(w, "Buy World")
	_ = w.Flush()
}
