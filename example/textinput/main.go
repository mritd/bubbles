package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/textinput"
)

func main() {
	m := textinput.Model{ValidateFunc: textinput.VFNotBlank}
	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}
}
