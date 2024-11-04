package utils

import (
	"os"

	"golang.org/x/term"
)

func GetTerminalSize() (width int, height int, err error) {
	return term.GetSize(int(os.Stdout.Fd()))
}
