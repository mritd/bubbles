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

// copy from textinput.Model
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

// Model is a data container used to store TUI status information,
// the ui rendering success style is as follows:
//
//	✔ Please Input: aaaa
type Model struct {
	// CharLimit is the maximum amount of characters this input element will
	// accept. If 0 or less, there's no limit.
	CharLimit int

	// Width is the maximum number of characters that can be displayed at once.
	// It essentially treats the text field like a horizontally scrolling
	// viewport. If 0 or less this setting is ignored.
	Width int

	// Prompt is the prefix of the prompt library, the user needs to define
	// the format(including spaces)
	Prompt string

	// PromptColor defines the color of the prompt prefix
	PromptColor string

	// ValidateFunc is a "real-time verification" function, which verifies
	// whether the terminal input data is legal in real time
	ValidateFunc func(string) error

	// ValidateOkPrefix is the prompt prefix when the validation fails
	ValidateOkPrefix string

	// ValidateErrPrefix is the prompt prefix when the verification is successful
	ValidateErrPrefix string

	// EchoMode sets the input behavior of the text input field.
	EchoMode EchoMode

	init    bool
	showErr bool
	input   textinput.Model
	err     error
}

// initData initialize the data model, set the default value and
// fix the wrong parameter settings during initialization
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

// Init performs some io initialization actions
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update method responds to various events and modifies the data model
// according to the corresponding events
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
			// When press the Enter button, if there is a verification error,
			// an error message is displayed.
			m.showErr = true
			if m.err == nil {
				return m, tea.Quit
			}
		case tea.KeyRunes:
			// Hide verification failure message when entering content again
			m.showErr = false
			m.err = nil
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		m.showErr = true
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	// Perform real-time verification function after each input
	m.err = m.ValidateFunc(m.input.Value())

	return m, cmd
}

// View reads the data state of the data model for rendering
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

// VFDoNothing is a verification function that does nothing
func VFDoNothing(_ string) error { return nil }

// VFNotBlank is a verification function that checks whether the input is empty
func VFNotBlank(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("input is empty")
	}
	return nil
}
