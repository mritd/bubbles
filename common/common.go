package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

const DONE = "DONE"

var term = termenv.ColorProfile()

// FontColor sets the color of the given string and bolds the font
func FontColor(str, color string) string {
	return termenv.String(str).Foreground(term.Color(color)).Bold().String()
}

// GenSpaces generate a space string of specified length
func GenSpaces(l int) string {
	return GenStr(l, " ")
}

// GenMask generate a mask string of the specified length
func GenMask(l int) string {
	return GenStr(l, "*")
}

// GenStr generate a string of the specified length, the string is composed of the given characters
func GenStr(l int, s string) string {
	var ss string
	for i := 0; i < l; i++ {
		ss += s
	}
	return ss
}

func Done() tea.Msg {
	return DONE
}
