package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/progressbar"
)

func main() {
	m := &progressbar.Model{
		Width:       40,
		InitMessage: "Initializing, please wait...",
		Stages: []progressbar.ProgressFunc{
			func() (string, error) {
				time.Sleep(time.Second)
				return "🐌 INFO: stage1", nil
			},
			func() (string, error) {
				time.Sleep(time.Second)
				return "🐌🐌 INFO: stage2", nil
			},
			func() (string, error) {
				time.Sleep(time.Second)
				return "🐌🐌🐌 INFO: stage3", nil
			},
			func() (string, error) {
				time.Sleep(time.Second)
				return "🐌🐌🐌🐌 INFO: stage4", nil
			},
			func() (string, error) {
				time.Sleep(time.Second)
				return "🐌🐌🐌🐌🐌 INFO: stage5", fmt.Errorf("🐞 Error: test error")
			},
		},
	}

	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}

	if m.Error() != nil {
		fmt.Printf("Stage func [%d] run failed: %s\n", m.Index()+1, m.Error())
	}
}
