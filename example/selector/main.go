package main

import (
	"fmt"
	"log"

	"github.com/mritd/bubbles/selector"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := &selector.Model{
		Data: []interface{}{
			"asdsdgsfgdsgsdfsd",
			"ijhlgmhldfgdg",
			"fgdfghdffgdgfdsgsdf",
			"opsdngdfgfhgkfgggdfg",
			"ghfgkhjgkksdfdsfs",
			"fgmhfkmhfghmjhkhjk",
			"chlfmhkmgfkjhfg",
			"kfgmhlfmjhgmhm",
			"lknkjas7yfgndndgldnflbah",
		},
		PerPage: 4,
		HeaderFunc: func(m selector.Model, prtIndex, drtIndex int) string {
			return fmt.Sprintf(selector.DefaultHeader+"\nCurrent Data index: %d\nCurrent Selected: %v", m.Index(), m.Selected())
		},
		SelectedFunc:   selector.DefaultSelectedFuncWithIndex("[%d]"),
		UnSelectedFunc: selector.DefaultUnSelectedFuncWithIndex(" %d."),
	}

	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}
	if !m.Canceled() {
		log.Printf("selected index => %d\n", m.Index())
		log.Printf("selected vaule => %s\n", m.Selected())
	} else {
		log.Println("user canceled...")
	}
}
