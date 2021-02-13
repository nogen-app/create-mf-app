package io

import (
	"fmt"

	"github.com/fatih/color"
)

type symbol string

const (
	//Checkmark is a checkmark symbol, that can be printed to the terminal
	Checkmark symbol = "✓"
	//Cross is a cross symbol, that can be printed to the terminal
	Cross symbol = "˟"
)

//PrintGreenCheckmark prints the msg prephased with a green checkmark to the stdout
//Used to indicate that someething was a success
func PrintGreenCheckmark(msg string) {
	green := color.New(color.FgGreen).PrintfFunc()
	green("\r[%s] ", Checkmark)
	fmt.Print(msg)
}

//PrintRedCross prints the msg prephased with a red cross to the stdout
//Used to indicate that a certain thing failed
func PrintRedCross(msg string) {
	red := color.New(color.FgRed).PrintfFunc()
	red("[%s] ", Cross)
	fmt.Print(msg)
}
