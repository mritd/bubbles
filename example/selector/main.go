package main

import (
	"fmt"
	"log"

	"github.com/mritd/bubbles/common"

	"github.com/mritd/bubbles/selector"

	tea "github.com/charmbracelet/bubbletea"
)

type CommitType string

const (
	FEAT     CommitType = "feat"
	FIX      CommitType = "fix"
	DOCS     CommitType = "docs"
	STYLE    CommitType = "style"
	REFACTOR CommitType = "refactor"
	TEST     CommitType = "test"
	CHORE    CommitType = "chore"
	PERF     CommitType = "perf"
)

type TypeMessage struct {
	Type          CommitType
	ZHDescription string
	ENDescription string
}

func main() {
	m := &selector.Model{
		Data: []interface{}{
			TypeMessage{Type: FEAT, ZHDescription: "新功能", ENDescription: "Introducing new features"},
			TypeMessage{Type: FIX, ZHDescription: "修复 Bug", ENDescription: "Bug fix"},
			TypeMessage{Type: DOCS, ZHDescription: "添加文档", ENDescription: "Writing docs"},
			TypeMessage{Type: STYLE, ZHDescription: "调整格式", ENDescription: "Improving structure/format of the code"},
			TypeMessage{Type: REFACTOR, ZHDescription: "重构代码", ENDescription: "Refactoring code"},
			TypeMessage{Type: TEST, ZHDescription: "增加测试", ENDescription: "When adding missing tests"},
			TypeMessage{Type: CHORE, ZHDescription: "CI/CD 变动", ENDescription: "Changing CI/CD"},
			TypeMessage{Type: PERF, ZHDescription: "性能优化", ENDescription: "Improving performance"},
		},
		PerPage: 5,
		// Use the arrow keys to navigate: ↓ ↑ → ←
		// Select Commit Type:
		HeaderFunc: selector.DefaultHeaderFuncWithAppend("Select Commit Type:"),
		// [1] feat (Introducing new features)
		SelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(TypeMessage)
			return common.FontColor(fmt.Sprintf("[%d] %s (%s)", gdIndex+1, t.Type, t.ENDescription), selector.ColorSelected)
		},
		// 2. fix (Bug fix)
		UnSelectedFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := obj.(TypeMessage)
			return common.FontColor(fmt.Sprintf(" %d. %s (%s)", gdIndex+1, t.Type, t.ENDescription), selector.ColorUnSelected)
		},
		// --------- Commit Type ----------
		// Type: feat
		// Description: 新功能(Introducing new features)
		FooterFunc: func(m selector.Model, obj interface{}, gdIndex int) string {
			t := m.PageSelected().(TypeMessage)
			footerTpl := `--------- Commit Type ----------
Type: %s
Description: %s(%s)`
			return common.FontColor(fmt.Sprintf(footerTpl, t.Type, t.ZHDescription, t.ENDescription), selector.ColorFooter)
		},
		FinishedFunc: func(s interface{}) string {
			return common.FontColor("Current selected: ", selector.ColorFinished) + string(s.(TypeMessage).Type) + "\n"
		},
	}

	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}
	if !m.Canceled() {
		log.Printf("selected index => %d\n", m.Index())
		log.Printf("selected vaule => %s\n", m.Selected())
		log.Printf("selected vaule => %s\n", m.PageSelected())
	} else {
		log.Println("user canceled...")
	}
}
