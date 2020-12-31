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

type ProgressFunc func() (string, error)

type Model struct {
	Width      int
	Stages     []ProgressFunc
	stageIndex int
	message    string
	err        error
	progress   float64
	loaded     bool
	init       bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	prompt := indent.String("\n"+m.message, 2)
	bar := indent.String("\n"+progressbar(m.Width, m.progress)+"%"+"\n\n", 2)
	return prompt + bar
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.init {
		m.initData()
		return m, func() tea.Msg {
			return m.Stages[m.stageIndex]
		}
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
		m.stageIndex++
		if !m.loaded {
			m.progress += float64(1) / float64(len(m.Stages))
			if m.progress >= 1+float64(1)/float64(len(m.Stages)) {
				m.progress = 1
				m.loaded = true
				return m, tea.Quit
			}
		}
	}

	return m, func() tea.Msg {
		return m.Stages[m.stageIndex-1]
	}
}

func (m *Model) initData() {
	m.stageIndex = 0
	if m.Width == 0 {
		m.Width = 35
	}
	m.init = true
}

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
