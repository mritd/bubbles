package main

import (
	"fmt"
	"log"

	"github.com/mritd/bubbles/common"

	"github.com/mritd/bubbles/selector"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	sl selector.Model
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// By default, the prompt component will not return a "tea.Quit"
	// message unless Ctrl+C is pressed.
	//
	// If there is no error in the input, the prompt component returns
	// a "common.DONE" message when the Enter key is pressed.
	switch msg {
	case common.DONE:
		return m, tea.Quit
	}

	_, cmd := m.sl.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.sl.View()
}

type TypeMessage struct {
	Type          string
	ZHDescription string
	ENDescription string
}

func main() {
	m := &model{
		sl: selector.Model{
			Data: []interface{}{
				TypeMessage{Type: "feat", ZHDescription: "新功能", ENDescription: "Introducing new features"},
				TypeMessage{Type: "fix", ZHDescription: "修复 Bug", ENDescription: "Bug fix"},
				TypeMessage{Type: "docs", ZHDescription: "添加文档", ENDescription: "Writing docs"},
				TypeMessage{Type: "style", ZHDescription: "调整格式", ENDescription: "Improving structure/format of the code"},
				TypeMessage{Type: "refactor", ZHDescription: "重构代码", ENDescription: "Refactoring code"},
				TypeMessage{Type: "test", ZHDescription: "增加测试", ENDescription: "When adding missing tests"},
				TypeMessage{Type: "chore", ZHDescription: "CI/CD 变动", ENDescription: "Changing CI/CD"},
				TypeMessage{Type: "perf", ZHDescription: "性能优化", ENDescription: "Improving performance"},
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
				t := m.Selected().(TypeMessage)
				footerTpl := `
Type: %s
Description: %s(%s)`
				return common.FontColor(fmt.Sprintf(footerTpl, t.Type, t.ZHDescription, t.ENDescription), selector.ColorFooter)
			},
			FinishedFunc: func(s interface{}) string {
				return common.FontColor("Current selected: ", selector.ColorFinished) + s.(TypeMessage).Type + "\n"
			},
		},
	}

	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}
	if !m.sl.Canceled() {
		log.Printf("selected index => %d\n", m.sl.Index())
		log.Printf("selected vaule => %s\n", m.sl.Selected())
	} else {
		log.Println("user canceled...")
	}
}
