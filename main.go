package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	home "chatip/home"
)

func main() {
	p := tea.NewProgram(home.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
