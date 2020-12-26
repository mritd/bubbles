package textinput

import (
	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	CharLimit         int
	Width             int
	Prompt            string
	ValidateOkPrefix  string
	ValidateErrPrefix string
	textInput         input.Model
	err               error
}

func (m *Model) initData() {
	m.textInput = input.NewModel()
	m.textInput.Focus()
	if m.CharLimit == 0 {
		m.CharLimit, m.textInput.CharLimit = 30, 30
	} else {
		m.textInput.CharLimit = m.CharLimit
	}
	if m.Width == 0 {
		m.Width, m.textInput.Width = 35, 35
	} else {
		m.textInput.Width = m.Width
	}
}

func (m Model) Init() tea.Cmd {
	return input.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyEsc:
			fallthrough
		case tea.KeyEnter:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {

	return m.textInput.View()
}
