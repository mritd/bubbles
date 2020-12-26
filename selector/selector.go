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
	// HeaderFunc 头部渲染函数
	HeaderFunc func(m Model, obj interface{}, gdIndex int) string
	// HeaderColor 头部渲染颜色
	HeaderColor string
	// Cursor 光标样式
	Cursor string
	// CursorColor 光标渲染颜色
	CursorColor string
	// SelectedFunc 被选中数据渲染函数
	SelectedFunc func(m Model, obj interface{}, gdIndex int) string
	// SelectedColor 被选中数据渲染颜色
	SelectedColor string
	// UnSelectedFunc 未选中数据渲染函数
	UnSelectedFunc func(m Model, obj interface{}, gdIndex int) string
	// UnSelectedColor 未选中数据渲染颜色
	UnSelectedColor string
	// FooterFunc 底部渲染函数
	FooterFunc func(m Model, obj interface{}, gdIndex int) string
	// FooterColor 底部渲染颜色
	FooterColor string
	// PerPage 每页数据量
	PerPage int
	// Data 要渲染的数据集合
	Data []interface{}

	// init 指示数据模型是否完成了初始化
	init bool
	// canceled 指示是否取消了操作
	canceled bool
	// pageData 当前页面实时渲染的数据集合
	pageData []interface{}
	// index 全局实时索引位置
	index int
	// maxIndex 全局允许的最大索引位置
	maxIndex int
	// pageIndex 当前页面实时索引位置
	pageIndex int
	// pageMaxIndex 当前页面允许的最大索引位置
	pageMaxIndex int
}

// Init 方法执行一些 I/O 初始化动作
func (m Model) Init() tea.Cmd {
	return nil
}

// View 方法读取数据模型数据状态进行渲染
func (m Model) View() string {
	// 光标直接渲染，无需处理
	cursor := common.FontColor(m.Cursor, m.CursorColor)
	// 模板函数可能会在列表头、尾和数据区动态显示，需要增加动态索引
	var header, data, footer string
	for i, obj := range m.pageData {
		// 光标前缀(选中行有光标，非选中行没有)
		var cursorPrefix string
		// 每行数据的渲染样式(选中行与非选中行颜色渲染等理论应当不同)
		var dataLine string
		// 全局动态索引 globalDynamicIndex 计算时考虑三种情况:
		//
		// 首页: pageIndex(页面实时索引)、index(全局实时索引) 两者保持一致
		//   1. feat (Introducing new features)
		//   2. fix (Bug fix)
		//   3. docs (Writing docs)
		//   4. style (Improving structure/format of the code)
		//   5. refactor (Refactoring code)
		//» [6] test (When adding missing tests)
		//
		// 滑动下翻页: pageIndex(页面实时索引) 固定到最大值，index(全局实时索引) 随滑动不断增大
		//   2. fix (Bug fix)
		//   3. docs (Writing docs)
		//   4. style (Improving structure/format of the code)
		//   5. refactor (Refactoring code)
		//   6. test (When adding missing tests)
		//» [7] chore (Changing CI/CD)
		//
		// 滑动上翻页: pageIndex(页面实时索引) 固定到最小值，index(全局实时索引) 随滑动不断减小
		//» [3] docs (Writing docs)
		//   4. style (Improving structure/format of the code)
		//   5. refactor (Refactoring code)
		//   6. test (When adding missing tests)
		//   7. chore (Changing CI/CD)
		//   8. perf (Improving performance)
		//
		// 三种情况下，m.index - m.pageIndex 的差值为 全局实时索引距离页面实时索引的距离，假定为 n
		// 那么遍历页面数据区时，遍历索引 i 理论上就是页面实时索引 pageIndex 的一个 "替身"
		// 此时 遍历索引 i + 全局实时索引距离页面实时索引的距离 n 即等于遍历时页面元素的全局索引位置 globalDynamicIndex
		globalDynamicIndex := i + (m.index - m.pageIndex)
		// 遍历数据区时，如果遍历的索引位置等于当前页面索引位置
		// 则当前数据区的数据为列表菜单选中数据，否则为未选中数据
		if i == m.pageIndex {
			// 光标与选中数据样式之间保持一个空格距离
			cursorPrefix = cursor + " "
			// m: 当前对象副本，将其传递到用户自定义渲染函数，方便用户读取一些状态信息，从而完成渲染
			// obj: 当前遍历到数据区的单条数据，将其传递到用户自定义渲染函数帮助用户得知当前需要渲染的数据
			// globalDynamicIndex: 当前遍历元素所对应的全局数据索引位置，将其传递到用户自定义渲染函数
			//                     帮助用户实现增加序号等渲染动作
			dataLine = common.FontColor(m.SelectedFunc(m, obj, globalDynamicIndex), m.SelectedColor) + "\n"
		} else {
			// 未选中行光标不显示，通过空白符对齐选中行
			cursorPrefix = common.GenSpaces(runewidth.StringWidth(m.Cursor) + 1)
			dataLine = common.FontColor(m.UnSelectedFunc(m, obj, globalDynamicIndex), m.UnSelectedColor) + "\n"
		}
		data += cursorPrefix + dataLine
		header = common.FontColor(m.HeaderFunc(m, obj, globalDynamicIndex), m.HeaderColor)
		footer = common.FontColor(m.FooterFunc(m, obj, globalDynamicIndex), m.FooterColor)
	}

	return fmt.Sprintf("%s\n\n%s\n%s\n", header, data, footer)
}

// Update 方法响应各种事件，根据对应事件修改数据模型
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

// Index 返回全局实时索引位置
func (m *Model) Index() int {
	return m.index
}

// PageIndex 返回页面实时索引位置
func (m *Model) PageIndex() int {
	return m.pageIndex
}

// PageData 返回当前页面数据区切片
func (m *Model) PageData() []interface{} {
	return m.pageData
}

// Selected 返回当前选中的元素
func (m *Model) Selected() interface{} {
	return m.Data[m.index]
}

// PageSelected 返回当前选中的元素(与 Selected 方法效果一致)
func (m *Model) PageSelected() interface{} {
	return m.pageData[m.pageIndex]
}

// Canceled 判断是否取消了操作
func (m *Model) Canceled() bool {
	return m.canceled
}

// MoveDown 方法执行向下移动光标动作，同时调整内部索引并刷新数据区
func (m *Model) MoveDown() {
	// 页面索引没有到达最大值，页面数据区不需要变更
	if m.pageIndex < m.pageMaxIndex {
		m.pageIndex++
		// 滑动前检测全局索引是否达到最大值
		if m.index < m.maxIndex {
			m.index++
		}
		return
	}

	// 页面索引到达最大值，滑动页面数据区窗口，页面索引维持最大值
	if m.pageIndex == m.pageMaxIndex {
		// 滑动前检测全局索引是否达到最大值
		if m.index < m.maxIndex {
			// 全局索引递增
			m.index++
			// 窗口向下滑动一个元素
			m.pageData = m.Data[m.index+1-m.PerPage : m.index+1]
			return
		}
	}
}

// MoveUp 方法执行向上移动光标动作，同时调整内部索引并刷新数据区
func (m *Model) MoveUp() {
	// 页面索引没有达到最小值，页面数据区不需要更新
	if m.pageIndex > 0 {
		m.pageIndex--
		// 滑动前检测全局索引是否达到最小值
		if m.index > 0 {
			m.index--
		}
		return
	}

	// 页面索引到达最小值，滑动页面数据窗口，页面索引维持最小值
	if m.pageIndex == 0 {
		// 滑动前检测全局索引是否达到最小值
		if m.index > 0 {
			// 窗口向上滑动一个元素
			m.pageData = m.Data[m.index-1 : m.index-1+m.PerPage]
			// 全局索引递减
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
