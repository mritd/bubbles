package progressbar

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

const (
	progressFullChar  = "█"
	progressEmptyChar = "░"
)

// General stuff for styling the view
var (
	term          = termenv.ColorProfile()
	subtle        = makeFgStyle("241")
	progressEmpty = subtle(progressEmptyChar)
)

// ProgressFunc is a simple function, the progress bar will step a certain distance after each execution
type ProgressFunc func() (string, error)

// Model is a data container used to store TUI status information.
type Model struct {
	Width       int
	Stages      []ProgressFunc
	InitMessage string
	stageIndex  int
	message     string
	err         error
	progress    float64
	loaded      bool
	init        bool
}

// Init performs some io initialization actions, The current Init returns the first ProgressFunc
// to trigger the program to run.
func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		if len(m.Stages) > 0 {
			return m.Stages[m.stageIndex]
		} else {
			return nil
		}
	}
}

// View reads the data state of the data model for rendering
func (m Model) View() string {
	prompt := indent.String("\n"+makeInfo(m.message), 2)
	if m.err != nil {
		prompt = indent.String("\n"+makeError(m.err.Error()), 2)
	}
	bar := indent.String("\n"+progressbar(m.Width, m.progress)+"%"+"\n\n", 2)
	return prompt + bar
}

// Update method responds to various events and modifies the data model
// according to the corresponding events
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.init {
		m.initData()
		m.message = makeInfo(m.InitMessage)
		return m, nil
	}

	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	if pf, ok := msg.(ProgressFunc); ok {
		m.message, m.err = pf()
		if m.err != nil {
			return m, tea.Quit
		}
		if m.stageIndex < len(m.Stages)-1 {
			m.stageIndex++
		}
		if !m.loaded {
			// The progress bar steps a certain distance after each successful execution
			m.progress += float64(1) / float64(len(m.Stages))
			// If all ProgressFunc has been executed, exit the TUI
			if m.progress > 1 {
				m.progress = 1
				m.loaded = true
				return m, tea.Quit
			}
		}
	}

	// Return to the next ProgressFunc, trigger the traversal
	return m, func() tea.Msg {
		return m.Stages[m.stageIndex]
	}
}

// initData initialize the data model, set the default value and
// fix the wrong parameter settings during initialization
func (m *Model) initData() {
	m.stageIndex = 0
	if m.Width == 0 {
		m.Width = 40
	}
	m.init = true
}

// Error returns the error generated during the execution of ProgressFunc
func (m *Model) Error() error {
	return m.err
}

// Index returns the index of the ProgressFunc currently executed
func (m *Model) Index() int {
	return m.stageIndex
}

// progressbar is responsible for rendering the progress bar UI
func progressbar(width int, percent float64) string {
	w := float64(width)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(progressFullChar).Foreground(term.Color(makeRamp("#B14FFF", "#00FFA3", w)[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Utils

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexidecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

func makeInfo(msg string) string {
	return termenv.String(msg).Foreground(term.Color("2")).Bold().String()
}

func makeError(msg string) string {
	return termenv.String(msg).Foreground(term.Color("9")).Bold().String()
}
