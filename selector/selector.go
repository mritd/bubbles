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
	maxIndex        int
	pageIndex       int
	pageMaxIndex    int
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

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "down":
			m.MoveDown()
		case "up":
			m.MoveUp()
			if m.pageIndex == 0 && m.index > 0 {
				m.pageData = m.Data[m.index-1 : m.index-1+(m.PerPage-1)]
			}

		case "pgdown", "right", "l", "k":
			m.NextPage()
		case "pgup", "left", "h", "j":
			m.PrePage()
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

func (m *Model) MoveDown() {
	//// 页面索引到达最大值
	//if m.pageIndex == (m.PerPage-1) && m.index < len(m.Data)-1 {
	//	m.pageData = m.Data[m.index-(m.PerPage-1) : m.index]
	//}

	// 页面索引没有到达最大值，页面数据区不需要变更
	if m.pageIndex < m.PerPage-1 {
		m.pageIndex++
		// 页面索引已到达最大值，再次增加将泄出
		// 此时保持索引最大值，滑动页面数据区窗口
	} else if m.pageIndex == m.PerPage-1 {
		// 滑动前检测全局索引是否达到最大值
		if m.index < len(m.Data)-1 {
			m.pageData = m.Data[m.index+1 : m.index+1+m.PerPage]
			m.index++
		}
	}
	if m.index < len(m.Data)-1 {
		m.index++
	}
}

func (m *Model) MoveUp() {
	if m.pageIndex > 0 {
		m.pageIndex--
	}
	if m.index > 0 {
		m.index--
	}
}

// NextPage 触发下翻页动作，翻页时忽略页面索引(pageIndex)位置
func (m *Model) NextPage() {
	// 全局索引没有到达最大值
	if m.index < len(m.Data)-1 {
		// m.PerPage-m.pageIndex = 当前页索引距离页尾差值
		// m.index+当前页索引距离页尾差值 = 当前页尾全局索引
		// 如果 当前页尾全局索引 < (全局索引总长度 - 每页长度)
		// 说明 当前页尾索引引距离最大值至少有一页长度剩余
		if m.index+((m.PerPage-1)-m.pageIndex) < (len(m.Data)-1)-(m.PerPage-1) {
			// 步进一页
			m.index += m.PerPage - 1
			m.pageData = m.Data[m.index : m.index+(m.PerPage-1)]
		} else {
			// 如果全局索引没有到达最大值且剩余小于一页
			// 直接步进到最大值
			m.index = len(m.Data) - 1
			m.pageData = m.Data[len(m.Data)-1-(m.PerPage-1):]
		}
	}
}

// PrePage 触发上翻页动作，翻页时忽略页面索引(pageIndex)位置
func (m *Model) PrePage() {
	// 全局索引没有到达最小值
	if m.index > 0 {
		// 全局索引距离最小值至少有一页长度剩余
		if len(m.Data)-1-m.index >= m.PerPage {
			// 步进一页
			m.pageData = m.Data[m.index-m.PerPage : m.index]
			m.index -= m.PerPage
		} else {
			// 如果全局索引没有到达最小值且剩余小于一页
			// 直接步进到最小值
			m.pageData = m.Data[:m.PerPage]
			m.index = 0
		}
	}
}

func (m *Model) initData() {
	m.pageIndex = 0
	m.pageMaxIndex = m.PerPage - 1
	m.pageData = m.Data[:m.PerPage]
	m.index = 0
	m.maxIndex = len(m.Data) - 1
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

func DefaultHeaderWithAppend(append string) func(m Model, prtIndex, drtIndex int) string {
	return func(m Model, prtIndex, drtIndex int) string {
		return DefaultHeader + "\n" + append
	}
}

func DefaultSelectedFuncWithIndex(indexFormat string) func(m Model, prtIndex, drtIndex int) string {
	return func(m Model, prtIndex, drtIndex int) string {
		return fmt.Sprintf(indexFormat+" %v", drtIndex+1, m.PageData()[prtIndex])
	}
}

func DefaultUnSelectedFuncWithIndex(indexFormat string) func(m Model, prtIndex, drtIndex int) string {
	return func(m Model, prtIndex, drtIndex int) string {
		return fmt.Sprintf(indexFormat+" %v", drtIndex+1, m.PageData()[prtIndex])
	}
}
