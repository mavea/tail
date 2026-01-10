package main

import "fmt"

func main() {
	_, _ = fmt.Print("Hello World", "\n")
	_, _ = fmt.Print("12345", "\n")
	_, _ = fmt.Print("67890", "\n")
	_, _ = fmt.Print("abcde", "\n")

	_, _ = fmt.Print("*\033[2Ag\033[2Ch\033[1Bi\033[3Dj")
	_, _ = fmt.Print("\033[2Fk\033[s\033[1El\033[2C\033[um")

	_, _ = fmt.Print("\n", "\n", "\n")
	_, _ = fmt.Print("\033[1;0Hf\033[1;3fe")
}

//Hello World
//km345
//lg89h
//abcjei
//*
