package main

import "fmt"

const (
	ANSIBlack int = iota
	ANSIRed
	ANSIGreen
	ANSIYellow
	ANSIBlue
	ANSIMagenta
	ANSICyan
	ANSIWhite
)

const clearCode = "\x1b[0m"

// ANSI escape code colors.
// Corresponds to the color enums.
var colors = [...]string{
	"\x1b[30m",
	"\x1b[31m",
	"\x1b[32m",
	"\x1b[33m",
	"\x1b[34m",
	"\x1b[35m",
	"\x1b[36m",
	"\x1b[37m",
}

func AnsiWrap(msg string, color int) string {
	return fmt.Sprintf("%s%s%s", colors[color], msg, clearCode)
}
