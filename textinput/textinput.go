package textinput

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mritd/bubbles/common"

	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	DefaultPrompt            = "Please Input: "
	DefaultValidateOkPrefix  = "✔"
	DefaultValidateErrPrefix = "✘"
)

type Model struct {
	CharLimit         int
	Width             int
	Prompt            string
	PromptColor       string
	ValidateFunc      func(string) error
	ValidateOkPrefix  string
	ValidateErrPrefix string

	init      bool
	textInput input.Model
	err       error
}

func (m *Model) initData() {
	m.textInput = input.NewModel()
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
	if m.PromptColor == "" {
		m.PromptColor = "2"
	}
	if m.Prompt == "" {
		m.Prompt, m.textInput.Prompt = common.FontColor(DefaultPrompt, "2"), common.FontColor(DefaultPrompt, "2")
	} else {
		m.textInput.Prompt = m.Prompt
	}
	if m.ValidateFunc == nil {
		m.ValidateFunc = VFDoNothing
	}
	if m.ValidateOkPrefix == "" {
		m.ValidateOkPrefix = DefaultValidateOkPrefix
	}
	if m.ValidateErrPrefix == "" {
		m.ValidateErrPrefix = DefaultValidateErrPrefix
	}
	m.textInput.Focus()
	m.init = true
}

func (m Model) Init() tea.Cmd {
	return input.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if !m.init {
		m.initData()
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.err = m.ValidateFunc(m.textInput.Value())
			if m.err == nil {
				return m, tea.Quit
			}
		case tea.KeyRunes: // 按键更新后清除终端验证提示
			m.err = nil
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
	var prefix, prompt, errMsg string
	if m.err != nil {
		prefix = common.FontColor(m.ValidateErrPrefix, "1")
		prompt = prefix + " " + m.textInput.View()
		errMsg = common.FontColor(fmt.Sprintf("%s ERROR: %s\n", m.ValidateErrPrefix, m.err.Error()), "1")
		return fmt.Sprintf("%s\n%s\n", prompt, errMsg)
	} else {
		prefix = common.FontColor(m.ValidateOkPrefix, "2")
		prompt = prefix + " " + m.textInput.View()
		return prompt + "\n"
	}
}

func VFDoNothing(_ string) error { return nil }
func VFNotBlank(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("input is empty")
	}
	return nil
}
