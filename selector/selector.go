package selector

import (
	"fmt"
	"strconv"

	"github.com/mritd/bubbles/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
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

// Use the arrow keys to navigate: ↓ ↑ → ←
// Select Commit Type:
//
// » [1] feat (Introducing new features)
//    2. fix (Bug fix)
//    3. docs (Writing docs)
//    4. style (Improving structure/format of the code)
//    5. refactor (Refactoring code)
//
// --------- Commit Type ----------
// Type: feat
// Description: 新功能(Introducing new features)
type Model struct {
	// HeaderFunc Header rendering function
	HeaderFunc func(m Model, obj interface{}, gdIndex int) string
	// HeaderColor header rendering color
	HeaderColor string
	// Cursor cursor rendering style
	Cursor string
	// CursorColor cursor rendering color
	CursorColor string
	// SelectedFunc selected data rendering function
	SelectedFunc func(m Model, obj interface{}, gdIndex int) string
	// SelectedColor selected data rendering color
	SelectedColor string
	// UnSelectedFunc unselected data rendering function
	UnSelectedFunc func(m Model, obj interface{}, gdIndex int) string
	// UnSelectedColor unselected data rendering color
	UnSelectedColor string
	// FooterFunc footer rendering function
	FooterFunc func(m Model, obj interface{}, gdIndex int) string
	// FooterColor footer rendering color
	FooterColor string
	// PerPage data count per page
	PerPage int
	// Data the data set to be rendered
	Data []interface{}

	// init indicates whether the data model has completed initialization
	init bool
	// canceled indicates whether the operation was cancelled
	canceled bool
	// pageData data set rendered in real time on the current page
	pageData []interface{}
	// index global real time index
	index int
	// maxIndex global max index
	maxIndex int
	// pageIndex real time index of current page
	pageIndex int
	// pageMaxIndex current page max index
	pageMaxIndex int
}

// Init init func performs some io initialization actions
func (m Model) Init() tea.Cmd {
	return nil
}

// View view func reads the data state of the data model for rendering
func (m Model) View() string {
	// the cursor only needs to be displayed correctly
	cursor := common.FontColor(m.Cursor, m.CursorColor)
	// template functions may be displayed dynamically at the head, tail and data area
	// of the list, and a dynamic index(globalDynamicIndex) needs to be added
	var header, data, footer string
	for i, obj := range m.pageData {
		// cursor prefix (selected lines need to be displayed,
		// non-selected lines need not be displayed)
		var cursorPrefix string
		// the rendering style of each row of data (the rendering color
		// of selected rows and non-selected rows is different)
		var dataLine string
		// consider three cases when calculating globalDynamicIndex:
		//
		// first page: pageIndex(real time page index)、index(global real time index) keep the two consistent
		//   1. feat (Introducing new features)
		//   2. fix (Bug fix)
		//   3. docs (Writing docs)
		//   4. style (Improving structure/format of the code)
		//   5. refactor (Refactoring code)
		//» [6] test (When adding missing tests)
		//
		// slide down to page: pageIndex fixed to maximum, index increasing with sliding
		//   2. fix (Bug fix)
		//   3. docs (Writing docs)
		//   4. style (Improving structure/format of the code)
		//   5. refactor (Refactoring code)
		//   6. test (When adding missing tests)
		//» [7] chore (Changing CI/CD)
		//
		// swipe up to page: pageIndex fixed to minimum, index decrease with sliding
		//» [3] docs (Writing docs)
		//   4. style (Improving structure/format of the code)
		//   5. refactor (Refactoring code)
		//   6. test (When adding missing tests)
		//   7. chore (Changing CI/CD)
		//   8. perf (Improving performance)
		//
		// in three cases, `m.index - m.pageIndex = n`, `n` is the distance between the global real-time
		// index and the page real-time index. when traversing the page data area, think of the traversal
		// index i as a real-time page index pageIndex, `i + n =` i corresponding global index
		globalDynamicIndex := i + (m.index - m.pageIndex)
		// when traversing the data area, if the traversed index is equal to the current page index,
		// the currently traversed data is the data selected in the list menu, otherwise it is unselected data
		if i == m.pageIndex {
			// keep a space between the cursor and the selected data style
			cursorPrefix = cursor + " "
			// m: A copy of the current object and pass it to the user-defined rendering function to facilitate
			//    the user to read some state information for rendering
			//
			// obj: The single data currently traversed to the data area; pass it to the user-defined rendering
			//      function to help users know the current data that needs to be rendered
			//
			// globalDynamicIndex: The global data index corresponding to the current traverse data; pass it
			//                     to the user-defined rendering function to help users achieve rendering
			//                     actions such as adding serial numbers
			dataLine = common.FontColor(m.SelectedFunc(m, obj, globalDynamicIndex), m.SelectedColor) + "\n"
		} else {
			// the cursor is not displayed on the unselected line, and the selected line is aligned with the blank character
			cursorPrefix = common.GenSpaces(runewidth.StringWidth(m.Cursor) + 1)
			dataLine = common.FontColor(m.UnSelectedFunc(m, obj, globalDynamicIndex), m.UnSelectedColor) + "\n"
		}
		data += cursorPrefix + dataLine
		header = common.FontColor(m.HeaderFunc(m, obj, globalDynamicIndex), m.HeaderColor)
		footer = common.FontColor(m.FooterFunc(m, obj, globalDynamicIndex), m.FooterColor)
	}

	return fmt.Sprintf("%s\n\n%s\n%s\n", header, data, footer)
}

// Update update method responds to various events and modifies the data model
// according to the corresponding events
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
		case "right", "pgdown", "l", "k":
			m.NextPage()
		case "left", "pgup", "h", "j":
			m.PrePage()
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.Forward(msg.String())
		}
	}
	return m, nil
}

// Index index return the global real time index
func (m *Model) Index() int {
	return m.index
}

// PageIndex pageIndex return the real time index of the page
func (m *Model) PageIndex() int {
	return m.pageIndex
}

// PageData pageData return the current page data area slice
func (m *Model) PageData() []interface{} {
	return m.pageData
}

// Selected selected return the currently selected data
func (m *Model) Selected() interface{} {
	return m.Data[m.index]
}

// PageSelected pageSelected return the currently selected data(same as the Selected func)
func (m *Model) PageSelected() interface{} {
	return m.pageData[m.pageIndex]
}

// Canceled canceled determine whether the operation is cancelled
func (m *Model) Canceled() bool {
	return m.canceled
}

// MoveDown moveDown executes the downward movement of the cursor,
// while adjusting the internal index and refreshing the data area
func (m *Model) MoveDown() {
	// the page index has not reached the maximum value, and the page
	// data area does not need to be updated
	if m.pageIndex < m.pageMaxIndex {
		m.pageIndex++
		// check whether the global index reaches the maximum value before sliding
		if m.index < m.maxIndex {
			m.index++
		}
		return
	}

	// the page index reaches the maximum value, slide the page data area window,
	// the page index maintains the maximum value
	if m.pageIndex == m.pageMaxIndex {
		// check whether the global index reaches the maximum value before sliding
		if m.index < m.maxIndex {
			// global index increment
			m.index++
			// window slide down one data
			m.pageData = m.Data[m.index+1-m.PerPage : m.index+1]
			return
		}
	}
}

// MoveUp moveUp performs an upward movement of the cursor,
// while adjusting the internal index and refreshing the data area
func (m *Model) MoveUp() {
	// the page index has not reached the minimum value, and the page
	// data area does not need to be updated
	if m.pageIndex > 0 {
		m.pageIndex--
		// check whether the global index reaches the minimum before sliding
		if m.index > 0 {
			m.index--
		}
		return
	}

	// the page index reaches the minimum value, slide the page data window,
	// and the page index maintains the minimum value
	if m.pageIndex == 0 {
		// check whether the global index reaches the minimum before sliding
		if m.index > 0 {
			// window slide up one data
			m.pageData = m.Data[m.index-1 : m.index-1+m.PerPage]
			// global index decrement
			m.index--
			return
		}
	}
}

// NextPage 触发下翻页动作，翻页时不变更页面实时索引(pageIndex)位置
func (m *Model) NextPage() {
	// 获取数据区对应的全局数据区切片起始和终止位置: m.Data[start:end]
	//
	// 注意: go 的切片是左闭右开的 [start,end)，假设全局数据区无限长度
	//      的情况下，页面数据区的 end 应当始终是实际页面长度+1，也就是说
	//      end 最大值在有限长度下应该等于 len(m.Data)
	pageStart, pageEnd := m.pageIndexInfo()
	// 数据区终止位置没有到达最大值时分为两种情况
	if pageEnd < len(m.Data) {
		// 数据区终止位置距离全局最大索引还有至少一个页面长度
		if len(m.Data)-pageEnd >= m.PerPage {
			// 页面数据区向后滑动一个页面
			m.pageData = m.Data[pageStart+m.PerPage : pageEnd+m.PerPage]
			// 直接将全局实时索引步进一个页面长度
			m.index += m.PerPage
		} else { // 全局实时索引距离全局最大索引小于一个页面长度
			// 页面数据区直接滑动到最后
			m.pageData = m.Data[len(m.Data)-m.PerPage : len(m.Data)]
			// 由于已经滑动到了最后，所以: 滑动距离 = 滑动最后的索引位置 - 未滑动前位置
			// 此时将全局实时索引也同步滑动距离即可
			m.index += len(m.Data) - pageEnd
		}
	}
}

// PrePage 触发上翻页动作，翻页时忽略页面索引(pageIndex)位置
func (m *Model) PrePage() {
	// 获取数据区对应的全局数据区切片起始和终止位置: m.Data[start:end]
	//
	// 注意: go 的切片是左闭右开的 [start,end)，假设全局数据区无限长度
	//      的情况下，页面数据区的 end 应当始终是实际页面长度+1，也就是说
	//      end 最大值在有限长度下应该等于 len(m.Data)
	pageStart, pageEnd := m.pageIndexInfo()
	// 数据区起始位置没有达到最小值时分为两种情况
	if pageStart > 0 {
		// 数据区起始位置距离最小值至少有一页长度剩余
		if pageStart >= m.PerPage {
			// 后退一页
			m.pageData = m.Data[pageStart-m.PerPage : pageEnd-m.PerPage]
			// 直接将全局实时索引后退一个页面长度
			m.index -= m.PerPage
		} else {
			// 如果数据区起始位置距离最小值小于一页
			// 直接后退到最小值
			m.pageData = m.Data[:m.PerPage]
			// 由于已经滑动到了最小值，所以: 滑动距离 = 滑动前位置 - 最小值0
			// 此时将全局实时索引也同步滑动距离即可
			m.index -= pageStart - 0
		}
	}
}

// Forward 触发快速跳转动作，如果按键不合法则维持原状
func (m *Model) Forward(pageIndex string) {
	// 输入层保证数据准确性，直接忽略 err
	idx, _ := strconv.Atoi(pageIndex)
	idx--

	// 目标索引位置已经超出页面最大索引，直接忽略
	if idx > m.pageMaxIndex {
		return
	}

	// 计算移动到目标长度
	l := idx - m.pageIndex
	// 全局索引移动
	m.index += l
	// 页面索引移动
	m.pageIndex = idx

}

// initData 负责初始化数据模型，初始化时会设置默认值以及修复错误的参数设置
func (m *Model) initData() {
	if m.PerPage > len(m.Data) || m.PerPage < 1 {
		m.PerPage = len(m.Data)
		m.pageData = m.Data
	} else {
		m.pageData = m.Data[:m.PerPage]
	}

	m.pageIndex = 0
	m.pageMaxIndex = m.PerPage - 1
	m.index = 0
	m.maxIndex = len(m.Data) - 1
	if m.HeaderFunc == nil {
		m.HeaderFunc = func(_ Model, _ interface{}, _ int) string { return DefaultHeader }
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
		m.SelectedFunc = func(m Model, obj interface{}, gdIndex int) string { return fmt.Sprint(obj) }
	}
	if m.SelectedColor == "" {
		m.SelectedColor = defaultSelectedColor
	}
	if m.UnSelectedFunc == nil {
		m.UnSelectedFunc = func(m Model, obj interface{}, gdIndex int) string { return fmt.Sprint(obj) }
	}
	if m.UnSelectedColor == "" {
		m.UnSelectedColor = defaultUnSelectedColor
	}
	if m.FooterFunc == nil {
		m.FooterFunc = func(_ Model, _ interface{}, _ int) string { return DefaultFooter }
	}
	if m.FooterColor == "" {
		m.FooterColor = defaultFooterColor
	}
	m.init = true
}

// pageIndexInfo 返回页面数据区对应全局数据区切片的起始、终止位置
func (m *Model) pageIndexInfo() (start, end int) {
	// 全局实时索引 - 全局页面实时索引 = 数据区起始索引
	start = m.index - m.pageIndex
	// 数据区起始位置 + 单页大小 = 数据区终止索引
	end = start + m.PerPage
	return
}

// DefaultHeaderFuncWithAppend 返回默认 HeaderFunc，并将给定字符串附加到默认头部下一行
func DefaultHeaderFuncWithAppend(append string) func(m Model, obj interface{}, gdIndex int) string {
	return func(m Model, obj interface{}, gdIndex int) string {
		return DefaultHeader + "\n" + append
	}
}

// DefaultSelectedFuncWithIndex 返回默认 SelectedFunc，并增加给定格式的序号前缀
func DefaultSelectedFuncWithIndex(indexFormat string) func(m Model, obj interface{}, gdIndex int) string {
	return func(m Model, obj interface{}, gdIndex int) string {
		return fmt.Sprintf(indexFormat+" %v", gdIndex+1, obj)
	}
}

// DefaultUnSelectedFuncWithIndex 返回默认 UnSelectedFunc，并增加给定格式的序号前缀
func DefaultUnSelectedFuncWithIndex(indexFormat string) func(m Model, obj interface{}, gdIndex int) string {
	return func(m Model, obj interface{}, gdIndex int) string {
		return fmt.Sprintf(indexFormat+" %v", gdIndex+1, obj)
	}
}
