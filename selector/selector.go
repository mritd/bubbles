package tplselector

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/termenv"
)

const (
	DefaultHeader          = "Use the arrow keys to navigate: ↓ ↑ → ←"
	DefaultFooter          = "Current page number details: %d/%d"
	DefaultCursor          = "»"
	defaultHeaderColor     = "15"
	defaultFooterColor     = "15"
	defaultCursorColor     = "2"
	defaultSelectedColor   = "14"
	defaultUnSelectedColor = "8"
)

var term = termenv.ColorProfile()

type Model struct {
	HeaderFunc      func(m Model, prtIndex, drtIndex int) string
	HeaderColor     string
	Cursor          string
	CursorColor     string
	SelectedFunc    func(m Model, prtIndex, drtIndex int) string
	SelectedColor   string
	UnSelectedFunc  func(m Model, prtIndex, drtIndex int) string
	UnSelectedColor string
	FooterFunc      func(m Model, prtIndex, drtIndex int) string
	FooterColor     string
	FooterShowIndex bool
	PerPage         int
	Data            []interface{}
	pageData        []interface{}
	init            bool
	canceled        bool
	index           int
	pageIndex       int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	cursor := fontColor(m.Cursor, m.CursorColor)

	var header, data, footer string
	for i := range m.pageData {
		var cursorPrefix string
		var dataLine string
		if i == m.pageIndex {
			cursorPrefix = cursor + " "
			dataLine = fontColor(m.SelectedFunc(m, i, i+m.index-m.pageIndex), m.SelectedColor) + "\n"
		} else {
			cursorPrefix = genSpaces(runewidth.StringWidth(m.Cursor) + 1)
			dataLine = fontColor(m.UnSelectedFunc(m, i, i+m.index-m.pageIndex), m.UnSelectedColor) + "\n"
		}
		data += cursorPrefix + dataLine
		header = fontColor(m.HeaderFunc(m, i, i+m.index-m.pageIndex), m.HeaderColor)
		if m.FooterShowIndex {
			footer = fontColor(fmt.Sprintf(m.FooterFunc(m, i, i+m.index-m.pageIndex), m.index+1, len(m.Data)), m.FooterColor)
		} else {
			footer = fontColor(m.FooterFunc(m, i, i+m.index-m.pageIndex), m.FooterColor)
		}
	}

	return fmt.Sprintf("%s\n\n%s\n%s\n", header, data, footer)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.init {
		m.initData()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.canceled = true
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		case "up", "pgup", "left", "h", "j":
			m.updatePageData()
			m.MoveUp()
		case "down", "pgdown", "right", "l", "k":
			m.MoveDown()
			m.updatePageData()
		}
	}
	return m, nil
}

func (m *Model) Index() int {
	return m.index
}

func (m *Model) PageIndex() int {
	return m.pageIndex
}

func (m *Model) PageData() []interface{} {
	return m.pageData
}

func (m *Model) Selected() interface{} {
	return m.Data[m.index]
}

func (m *Model) PageSelected() interface{} {
	return m.pageData[m.pageIndex]
}

func (m *Model) Canceled() bool {
	return m.canceled
}

func (m *Model) MoveUp() {
	if m.pageIndex > 0 {
		m.pageIndex--
	}
	if m.index > 0 {
		m.index--
	}
}

func (m *Model) MoveDown() {
	if m.pageIndex < m.PerPage-1 {
		m.pageIndex++
	}
	if m.index < len(m.Data)-1 {
		m.index++
	}
}

func (m *Model) initData() {
	m.pageData = m.Data[:m.PerPage]
	if m.HeaderFunc == nil {
		m.HeaderFunc = func(_ Model, _, _ int) string { return DefaultHeader }
	}
	if m.HeaderColor == "" {
		m.HeaderColor = defaultHeaderColor
	}
	if m.Cursor == "" {
		m.Cursor = DefaultCursor
	}
	if m.CursorColor == "" {
		m.CursorColor = defaultCursorColor
	}
	if m.SelectedFunc == nil {
		m.SelectedFunc = func(m Model, prtIndex, drtIndex int) string { return fmt.Sprint(m.pageData[prtIndex]) }
	}
	if m.SelectedColor == "" {
		m.SelectedColor = defaultSelectedColor
	}
	if m.UnSelectedFunc == nil {
		m.UnSelectedFunc = func(m Model, prtIndex, drtIndex int) string { return fmt.Sprint(m.pageData[prtIndex]) }
	}
	if m.UnSelectedColor == "" {
		m.UnSelectedColor = defaultUnSelectedColor
	}
	if m.FooterFunc == nil {
		m.FooterFunc = func(_ Model, _, _ int) string { return DefaultFooter }
		m.FooterShowIndex = true
	}
	if m.FooterColor == "" {
		m.FooterColor = defaultFooterColor
	}
	m.init = true
}

func (m *Model) updatePageData() {
	// 达到最大值，向下滑动窗口
	if m.pageIndex == m.PerPage-1 && m.index < len(m.Data) {
		m.pageData = m.Data[m.index+1-m.PerPage : m.index+1]
	}
	// 达到最小值，向上滑动窗口
	if m.pageIndex == 0 && m.index > 0 {
		m.pageData = m.Data[m.index-1 : m.index-1+m.PerPage]
	}
}

func fontColor(str, color string) string {
	return termenv.String(str).Foreground(term.Color(color)).Bold().String()
}

func genSpaces(l int) string {
	var s string
	for i := 0; i < l; i++ {
		s += " "
	}
	return s
}
