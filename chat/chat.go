package chat

import (
	// golang imports
	"fmt"
	"os"
	"os/exec"

	// github imports
	tea "github.com/charmbracelet/bubbletea"
)

type msg_chat struct {
	who		  string
	message   string
	time      string
	date      string
	send      bool
	delivered bool
}

type chat struct {
	ip        string
	own_name  string
	they_name []string
	msg  	  []msg_chat
	cursor    int
}

func InitialModel(ip string, ownName string) chat {
	return chat{
		ip:       ip,
		own_name: ownName,
	}
}

func (c chat) Init() tea.Cmd {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}

func (c chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit

		case "ctrl+m":
			// enter in message mode
			return c, nil

		case "esc":
			// get out of message mode
			return c, nil

		case "enter":
			// send message if in message mode
			return c, nil

		case "ctrl+e":
			// edit mode
			return c, nil

		case "ctrl+s":
			// export chat
			return c, nil
		}
	}

	return c, nil
}

func (c chat) View() string {
	// chat mode
	s := fmt.Sprintf("%s %s", c.ip, c.own_name)

	return s
}

func main() {
	p := tea.NewProgram(InitialModel("0.0.0.0", "default"))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
