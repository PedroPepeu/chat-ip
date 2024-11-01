package home

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"

	chat "chatip/chat"
)

type homeScreen struct {
	ip       string
	own_name string
	cursor 	 int
}

func InitialModel() homeScreen {
	return homeScreen{
		ip:       "127.0.0.1",
		own_name: "default",
	}
}

func (hs homeScreen) Init() tea.Cmd {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}

func (hs homeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return hs, tea.Quit

		case "up", "k":
			if hs.cursor > 0 {
				hs.cursor--
			}

		case "down", "j":
			// if hs.cursor < len(hs.choices)-1 {
			// 	hs.cursor++
			// }

		case "enter", " ":
			return chat.InitialModel().Update(msg)
			
		}
	}
	return hs, nil
}

func (hs homeScreen) View() string {
	s := "Home screen"
	return s
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
