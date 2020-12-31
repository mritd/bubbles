package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/progressbar"
)

func main() {
	m := progressbar.Model{
		Width: 35,
		Stages: []progressbar.ProgressFunc{
			func() (string, error) {
				return "stage1", nil
			},
			func() (string, error) {
				return "stage2", nil
			},
			func() (string, error) {
				return "stage3", nil
			},
			func() (string, error) {
				return "stage4", nil
			},
		},
	}

	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}
}
