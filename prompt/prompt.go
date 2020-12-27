package prompt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mritd/bubbles/common"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	DefaultPrompt            = "Please Input: "
	DefaultValidateOkPrefix  = "✔"
	DefaultValidateErrPrefix = "✘"
	defaultPromptOkColor     = "2"
	defaultValidateOkColor   = "2"
	defaultValidateErrColor  = "1"
)

// EchoMode sets the input behavior of the text input field.
type EchoMode int

const (
	// EchoNormal displays text as is. This is the default behavior.
	EchoNormal EchoMode = iota

	// EchoPassword displays the EchoCharacter mask instead of actual
	// characters.  This is commonly used for password fields.
	EchoPassword

	// EchoNone displays nothing as characters are entered. This is commonly
	// seen for password fields on the command line.
	EchoNone
)

type Model struct {
	CharLimit         int
	Width             int
	Prompt            string
	PromptColor       string
	ValidateFunc      func(string) error
	ValidateOkPrefix  string
	ValidateErrPrefix string
	EchoMode          int

	init    bool
	showErr bool
	input   textinput.Model
	err     error
}

func (m *Model) initData() {
	m.input = textinput.NewModel()
	if m.CharLimit == 0 {
		m.CharLimit, m.input.CharLimit = 30, 30
	} else {
		m.input.CharLimit = m.CharLimit
	}
	if m.Width == 0 {
		m.Width, m.input.Width = 35, 35
	} else {
		m.input.Width = m.Width
	}
	if m.PromptColor == "" {
		m.PromptColor = defaultPromptOkColor
	}
	if m.Prompt == "" {
		m.Prompt, m.input.Prompt = common.FontColor(DefaultPrompt, m.PromptColor), common.FontColor(DefaultPrompt, m.PromptColor)
	} else {
		m.input.Prompt = common.FontColor(m.Prompt, m.PromptColor)
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
	m.input.EchoMode = textinput.EchoMode(m.EchoMode)
	m.input.Focus()
	m.init = true
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.showErr = true
			if m.err == nil {
				return m, tea.Quit
			}
		case tea.KeyRunes:
			m.showErr = false
			m.err = nil
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	m.err = m.ValidateFunc(m.input.Value())

	return m, cmd
}

func (m Model) View() string {
	var prefix, prompt, errMsg string
	if m.err != nil {
		prefix = common.FontColor(m.ValidateErrPrefix, defaultValidateErrColor)
		prompt = prefix + " " + m.input.View()
		if m.showErr {
			errMsg = common.FontColor(fmt.Sprintf("%s ERROR: %s\n", m.ValidateErrPrefix, m.err.Error()), defaultValidateErrColor)
			return fmt.Sprintf("%s\n%s\n", prompt, errMsg)
		}
	} else {
		prefix = common.FontColor(m.ValidateOkPrefix, defaultValidateOkColor)
		prompt = prefix + " " + m.input.View()
	}

	return prompt + "\n"
}

func (m *Model) Value() string {
	return m.input.Value()
}

func VFDoNothing(_ string) error { return nil }

func VFNotBlank(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("input is empty")
	}
	return nil
}
