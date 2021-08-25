package main

import (
	"github.com/mritd/bubbles/common"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/prompt"
)

type model struct {
	input *prompt.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// By default, the prompt component will not return a "tea.Quit"
	// message unless Ctrl+C is pressed.
	//
	// If there is no error in the input, the prompt component returns
	// a "common.DONE" message when the Enter key is pressed.
	switch msg {
	case common.DONE:
		return m, tea.Quit
	}

	_, cmd := m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.input.View()
}

func (m model) Value() string {
	return m.input.Value()
}

func main() {
	m := model{input: &prompt.Model{ValidateFunc: prompt.VFNotBlank}}
	p := tea.NewProgram(&m)
	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(m.Value())
}
