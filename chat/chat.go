package chat

import (
	// golang imports
	"fmt"
	"os"
	"os/exec"
	"strings"

	// github imports
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type msg_chat struct {
	who       string
	message   string
	time      string
	send      bool
	delivered bool
}

type (
	errMsg error
)

type chat struct {
	viewport	viewport.Model
	ip          string
	own_name    string
	they_name   []string
	msg         []msg_chat
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func InitialModel(ip string, ownName string) chat {
	ta := textarea.New()
	ta.Placeholder = fmt.Sprintf("Connected to %s, as %s", ip, ownName)
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent("Type a message and press (enter) to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return chat{
		ip:			 ip,
		own_name:	 ownName,
		textarea:    ta,
		msg:         []msg_chat{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (c chat) Init() tea.Cmd {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	return textarea.Blink
}

func (c chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	c.textarea, tiCmd = c.textarea.Update(msg)
	c.viewport, vpCmd = c.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		// close chat
		case tea.KeyCtrlC, tea.KeyEsc:
			return c, tea.Quit

		// send message
		case tea.KeyEnter:

			// type msg_chat struct {
			// 	who       string
			// 	message   string
			// 	time      string
			// 	send      bool
			// 	delivered bool
			// }

			message := msg_chat {
				who:		c.own_name,
				message:	c.msg,
				time:		time.Now(),
				send:		true,
				delivered:	true,
			}

			c.msg = append(message, c.senderStyle.Render(c.own_name)+c.textarea.Value())
			c.viewport.SetContent(strings.Join(c.msg, "\n"))
			c.textarea.Reset()
			c.viewport.GotoBottom()

		// handling errors
		case errMsg:
			c.err = msg
			return c, nil

		// export chat
		case tea.KeyCtrlS:
			return c, nil
		}
	}

	return c, tea.Batch(tiCmd, vpCmd)
}

func (c chat) View() string {
	// chat mode
	s := fmt.Sprintf("%s\n\n%s",
		c.viewport.View(),
		c.textarea.View(),
	) + "\n\n"

	return s
}

func main() {
	p := tea.NewProgram(InitialModel("0.0.0.0", "default"))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
