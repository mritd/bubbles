package selector

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
	Header          string
	HeaderColor     string
	Cursor          string
	CursorColor     string
	SelectedColor   string
	UnSelectedColor string
	Footer          string
	FooterColor     string
	FooterShowIndex bool
	PerPage         int
	Data            []string
	init            bool
	index           int
	pageIndex       int
	pageData        []string
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	header := fontColor(m.Header, m.HeaderColor)
	cursor := fontColor(m.Cursor, m.CursorColor)

	var footerTpl string
	if m.FooterShowIndex {
		footerTpl = fmt.Sprintf(m.Footer, m.index+1, len(m.Data))
	} else {
		footerTpl = m.Footer
	}
	footer := fontColor(footerTpl, m.FooterColor)

	var data string
	for i, field := range m.pageData {
		var cursorPrefix string
		var dataLine string
		if i == m.pageIndex {
			cursorPrefix = cursor + " "
			dataLine = fontColor(field, m.SelectedColor) + "\n"
		} else {
			cursorPrefix = genSpaces(runewidth.StringWidth(m.Cursor) + 1)
			dataLine = fontColor(field, m.UnSelectedColor) + "\n"
		}
		data += cursorPrefix + dataLine
	}

	return fmt.Sprintf("%s\n\n%s\n%s\n", header, data, footer)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.init {
		m.initData()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.index = -1
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

func (m *Model) Selected() int {
	return m.index
}

func (m *Model) Canceled() bool {
	return m.index == -1
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
	if m.Header == "" {
		m.Header = DefaultHeader
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
	if m.SelectedColor == "" {
		m.SelectedColor = defaultSelectedColor
	}
	if m.UnSelectedColor == "" {
		m.UnSelectedColor = defaultUnSelectedColor
	}
	if m.Footer == "" {
		m.Footer = DefaultFooter
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
