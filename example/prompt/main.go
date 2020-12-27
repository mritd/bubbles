package main

import (
	"log"

	"github.com/mritd/bubbles/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/prompt"
)

func main() {
	m := &prompt.Model{ValidateFunc: prompt.VFNotBlank}
	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(m.Value())

	m1 := &prompt.Model{
		ValidateFunc: prompt.VFNotBlank,
		Prompt:       common.FontColor("Please input password: ", "2"),
		EchoMode:     prompt.EchoPassword,
	}
	p1 := tea.NewProgram(m1)
	err1 := p1.Start()
	if err1 != nil {
		log.Fatal(err1)
	}
	log.Println(m1.Value())
}
