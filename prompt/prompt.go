package prompt

import (
	"errors"
	"fmt"
	"github.com/mritd/bubbles/common"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	DefaultPrompt            = "Please Input: "
	DefaultValidateOkPrefix  = "✔"
	DefaultValidateErrPrefix = "✘"

	ColorPrompt      = "2"
	colorValidateOk  = "2"
	colorValidateErr = "1"
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

	// ValidateFunc is a "real-time verification" function, which verifies
	// whether the terminal input data is legal in real time
	ValidateFunc func(string) error

	// ValidateOkPrefix is the prompt prefix when the validation fails
	ValidateOkPrefix string

	// ValidateErrPrefix is the prompt prefix when the verification is successful
	ValidateErrPrefix string

	// EchoMode sets the input behavior of the text input field.
	EchoMode EchoMode

	init     bool
	canceled bool
	finished bool
	showErr  bool
	err      error

	input textinput.Model
}

// initData initialize the data model, set the default value and
// fix the wrong parameter settings during initialization
func (m *Model) initData() {
	if m.ValidateFunc == nil {
		m.ValidateFunc = VFDoNothing
	}
	if m.ValidateOkPrefix == "" {
		m.ValidateOkPrefix = DefaultValidateOkPrefix
	}
	if m.ValidateErrPrefix == "" {
		m.ValidateErrPrefix = DefaultValidateErrPrefix
	}
	if m.Prompt == "" {
		m.Prompt = common.FontColor(DefaultPrompt, ColorPrompt)
	}

	in := textinput.NewModel()
	in.CharLimit = m.CharLimit
	in.Width = m.Width
	in.Prompt = m.Prompt
	in.EchoMode = textinput.EchoMode(m.EchoMode)
	in.Focus()

	m.input = in
	m.init = true
}

// View reads the data state of the data model for rendering
func (m Model) View() string {
	if m.finished {
		switch m.EchoMode {
		case EchoNormal:
			return common.FontColor(m.ValidateOkPrefix, colorValidateOk) + " " + m.Prompt + m.Value() + "\n"
		case EchoNone:
			return common.FontColor(m.ValidateOkPrefix, colorValidateOk) + " " + m.Prompt + "\n"
		case EchoPassword:
			return common.FontColor(m.ValidateOkPrefix, colorValidateOk) + " " + m.Prompt + common.GenMask(len([]rune(m.Value()))) + "\n"
		}
	}

	var prompt, errMsg string
	if m.err != nil {
		prompt = common.FontColor(m.ValidateErrPrefix, colorValidateErr) + " " + m.input.View()
		if m.showErr {
			errMsg = common.FontColor(fmt.Sprintf("%s ERROR: %s\n", m.ValidateErrPrefix, m.err.Error()), colorValidateErr)
			return fmt.Sprintf("%s\n%s\n", prompt, errMsg)
		}
	} else {
		prompt = common.FontColor(m.ValidateOkPrefix, colorValidateOk) + " " + m.input.View()
	}

	return prompt + "\n"
}

// Update method responds to various events and modifies the data model
// according to the corresponding events
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	if !m.init {
		m.initData()
		return m, nil
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// We intercept some key events, because we need to handle it in the upper layer
		switch msg.Type {
		case tea.KeyCtrlC:
			// Terminate the UI program when Ctrl+C is pressed
			m.canceled = true
			return m, tea.Quit
		case tea.KeyEnter:
			// If the real-time verification function does not return an error,
			// then the input has been completed
			if m.err == nil {
				m.finished = true
				return m, common.Done
			}

			// If there is a verification error, the error message should be display
			m.showErr = true
		case tea.KeyRunes:
			// Hide verification failure message when entering content again
			m.showErr = false
			m.err = nil
		}

		// Call the underlying textinput to update the terminal display
		m.input, cmd = m.input.Update(msg)
		// Perform real-time verification function after each input
		m.err = m.ValidateFunc(m.input.Value())

	// We handle errors just like any other message
	// Note: msg is error only when there is an unexpected error in the underlying textinput
	case error:
		m.err = msg
		m.showErr = true
		return m, nil
	}

	return m, cmd
}

// Value return the input string
func (m Model) Value() string {
	return m.input.Value()
}

// Canceled determine whether the operation is cancelled
func (m Model) Canceled() bool {
	return m.canceled
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
