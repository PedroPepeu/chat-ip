package home

import (
	// golang imports
	"fmt"
	"os"
	"os/exec"
	"strings"

	// github imports
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// project imports
	chat "chatip/chat"
	utils "chatip/utils"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type homeScreen struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func InitialModel() homeScreen {
	hs := homeScreen{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range hs.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 15

		switch i {
		case 0:
			t.Placeholder = "0.0.0.0"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle

		case 1:
			t.Placeholder = "Alias"
			t.CharLimit = 20
		}

		hs.inputs[i] = t
	}

	return hs
}

func (hs homeScreen) Init() tea.Cmd {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	return textinput.Blink
}

func (hs homeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		// exit the application
		case "ctrl+c", "esc":
			return hs, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			hs.cursorMode++
			if hs.cursorMode > cursor.CursorHide {
				hs.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(hs.inputs))
			for i := range hs.inputs {
				cmds[i] = hs.inputs[i].Cursor.SetMode(hs.cursorMode)
			}
			return hs, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && hs.focusIndex == len(hs.inputs) {
				ip := hs.inputs[0].Value()
				alias := hs.inputs[1].Value()
				return chat.InitialModel(ip, alias).Update(msg)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				hs.focusIndex--
			} else {
				hs.focusIndex++
			}

			if hs.focusIndex > len(hs.inputs) {
				hs.focusIndex = 0
			} else if hs.focusIndex < 0 {
				hs.focusIndex = len(hs.inputs)
			}

			cmds := make([]tea.Cmd, len(hs.inputs))
			for i := 0; i <= len(hs.inputs)-1; i++ {
				if i == hs.focusIndex {
					// Set focused state
					cmds[i] = hs.inputs[i].Focus()
					hs.inputs[i].PromptStyle = focusedStyle
					hs.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				hs.inputs[i].Blur()
				hs.inputs[i].PromptStyle = noStyle
				hs.inputs[i].TextStyle = noStyle
			}

			return hs, tea.Batch(cmds...)
		}
	}

	cmd := hs.updateInputs(msg)
	return hs, cmd
}

func (hs *homeScreen) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(hs.inputs))

	for i := range hs.inputs {
		hs.inputs[i], cmds[i] = hs.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (hs homeScreen) View() string {
	var b strings.Builder

	for i := range hs.inputs {
		b.WriteString(hs.inputs[i].View())
		if i < len(hs.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if hs.focusIndex == len(hs.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	width, height, err := utils.GetTerminalSize()

	if err != nil {
		return "Error retrieving terminal size."
	}

	centeredStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(width).
		PaddingTop(height / 2)

	return centeredStyle.Render(b.String())
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
