package main

import (
	"log"

	"github.com/mritd/bubbles/selector"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := &selector.Model{
		Data: []string{
			"1 asdsdgsfgdsgsdfsd",
			"2 sadsdgdfghdfhfd",
			"3 fgdfghdffgdgfdsgsdf",
			"4 fgfghklfghkfnhgf",
			"5 ghfgkhjgkksdfdsfs",
			"6 fgmhfkmhfghmjhkhjk",
			"7 fhlfmhkmgfkjhfg",
			"8 kfgmhlfmjhgmhm",
			"9 lknkjas7yfgndndgldnflbah",
		},
		PerPage: 4,
		Header:  selector.DefaultHeader + "\nSelect Login Server:",
	}

	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}
	if !m.Canceled() {
		log.Printf("selected index => %d\n", m.Selected())
		log.Printf("selected vaule => %s\n", m.Data[m.Selected()])
	} else {
		log.Println("user canceled...")
	}
}
