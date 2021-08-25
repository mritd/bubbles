// Package selector is a terminal single-selection list library. selector library provides the
// functions of page up and down and key movement, and supports custom rendering methods.
package selector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mritd/bubbles/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
)

const (
	DefaultHeader   = "Use the arrow keys to navigate: ↓ ↑ → ←"
	DefaultFooter   = "Current page number details: %d/%d"
	DefaultCursor   = "»"
	DefaultFinished = "Current selected: %s\n"

	ColorHeader     = "15"
	ColorFooter     = "15"
	ColorCursor     = "2"
	ColorFinished   = "2"
	ColorSelected   = "14"
	ColorUnSelected = "8"
)

// Model is a data container used to store TUI status information,
// the ui rendering success style is as follows:
//
//	Use the arrow keys to navigate: ↓ ↑ → ←
// 	Select Commit Type:
//
// 	» [1] feat (Introducing new features)
//	   2. fix (Bug fix)
//	   3. docs (Writing docs)
//	   4. style (Improving structure/format of the code)
//	   5. refactor (Refactoring code)
//
//	--------- Commit Type ----------
//	Type: feat
//	Description: 新功能(Introducing new features)
type Model struct {
	// HeaderFunc Header rendering function
	HeaderFunc func(m Model, obj interface{}, gdIndex int) string
	// Cursor cursor rendering style
	Cursor string
	// CursorColor cursor rendering color
	CursorColor string
	// SelectedFunc selected data rendering function
	SelectedFunc func(m Model, obj interface{}, gdIndex int) string
	// UnSelectedFunc unselected data rendering function
	UnSelectedFunc func(m Model, obj interface{}, gdIndex int) string
	// FooterFunc footer rendering function
	FooterFunc func(m Model, obj interface{}, gdIndex int) string
	// FinishedFunc finished rendering function
	FinishedFunc func(selected interface{}) string
	// PerPage data count per page
	PerPage int
	// Data the data set to be rendered
	Data []interface{}

	// init indicates whether the data model has completed initialization
	init bool
	// canceled indicates whether the operation was cancelled
	canceled bool
	// finished indicates whether it has exited
	finished bool
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

// View reads the data state of the data model for rendering
func (m Model) View() string {
	if m.finished {
		return m.FinishedFunc(m.Selected())
	}

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
		//	  1. feat (Introducing new features)
		//	  2. fix (Bug fix)
		//	  3. docs (Writing docs)
		//	  4. style (Improving structure/format of the code)
		//	  5. refactor (Refactoring code)
		//	» [6] test (When adding missing tests)
		//
		// slide down to page: pageIndex fixed to maximum, index increasing with sliding
		//	  2. fix (Bug fix)
		//	  3. docs (Writing docs)
		//	  4. style (Improving structure/format of the code)
		//	  5. refactor (Refactoring code)
		//	  6. test (When adding missing tests)
		//	» [7] chore (Changing CI/CD)
		//
		// swipe up to page: pageIndex fixed to minimum, index decrease with sliding
		//	» [3] docs (Writing docs)
		//	  4. style (Improving structure/format of the code)
		//	  5. refactor (Refactoring code)
		//	  6. test (When adding missing tests)
		//	  7. chore (Changing CI/CD)
		//	  8. perf (Improving performance)
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
			dataLine = m.SelectedFunc(m, obj, globalDynamicIndex) + "\n"
		} else {
			// the cursor is not displayed on the unselected line, and the selected line is aligned with the blank character
			cursorPrefix = common.GenSpaces(runewidth.StringWidth(m.Cursor) + 1)
			dataLine = m.UnSelectedFunc(m, obj, globalDynamicIndex) + "\n"
		}
		data += cursorPrefix + dataLine
		header = m.HeaderFunc(m, obj, globalDynamicIndex)
		footer = m.FooterFunc(m, obj, globalDynamicIndex)
	}

	return fmt.Sprintf("%s\n\n%s\n%s", header, data, footer)
}

// Update method responds to various events and modifies the data model
// according to the corresponding events
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	if !m.init {
		m.initData()
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "q", "ctrl+c":
			m.canceled = true
			return m, tea.Quit
		case "enter":
			m.finished = true
			return m, common.Done
		case "down":
			m.moveDown()
		case "up":
			m.moveUp()
		case "right", "pgdown", "l", "k":
			m.nextPage()
		case "left", "pgup", "h", "j":
			m.prePage()
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.forward(msg.String())
		}
	}
	return m, nil
}

// moveDown executes the downward movement of the cursor,
// while adjusting the internal index and refreshing the data area
func (m *Model) moveDown() {
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

// moveUp performs an upward movement of the cursor,
// while adjusting the internal index and refreshing the data area
func (m *Model) moveUp() {
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

// nextPage triggers the page-down action, and does not change
// the real-time page index(pageIndex)
func (m *Model) nextPage() {
	// Get the start and end position of the page data area slice: m.Data[start:end]
	//
	// note: the slice is closed left and opened right: `[start,end)`
	//       assuming that the global data area has unlimited length,
	//       end should always be the actual page `length+1`,
	//       the maximum value of end should be equal to `len(m.Data)`
	//       under limited length
	pageStart, pageEnd := m.pageIndexInfo()
	// there are two cases when `end` does not reach the maximum value
	if pageEnd < len(m.Data) {
		// the `end` value is at least one page length away from the global maximum index
		if len(m.Data)-pageEnd >= m.PerPage {
			// slide back one page in the page data area
			m.pageData = m.Data[pageStart+m.PerPage : pageEnd+m.PerPage]
			// Global real-time index increases by one page length
			m.index += m.PerPage
		} else { // `end` is less than a page length from the global maximum index
			// slide the page data area directly to the end
			m.pageData = m.Data[len(m.Data)-m.PerPage : len(m.Data)]
			// `sliding distance` = `position after sliding` - `position before sliding`
			// the global real-time index should also synchronize the same sliding distance
			m.index += len(m.Data) - pageEnd
		}
	}
}

// prePage triggers the page-up action, and does not change
// the real-time page index(pageIndex)
func (m *Model) prePage() {
	// Get the start and end position of the page data area slice: m.Data[start:end]
	//
	// note: the slice is closed left and opened right: `[start,end)`
	//       assuming that the global data area has unlimited length,
	//       end should always be the actual page `length+1`,
	//       the maximum value of end should be equal to `len(m.Data)`
	//       under limited length
	pageStart, pageEnd := m.pageIndexInfo()
	// there are two cases when `start` does not reach the minimum value
	if pageStart > 0 {
		// `start` is at least one page length from the minimum
		if pageStart >= m.PerPage {
			// slide the page data area forward one page
			m.pageData = m.Data[pageStart-m.PerPage : pageEnd-m.PerPage]
			// Global real-time index reduces the length of one page
			m.index -= m.PerPage
		} else { // `start` to the minimum value less than one page length
			// slide the page data area directly to the start
			m.pageData = m.Data[:m.PerPage]
			// `sliding distance` = `position before sliding` - `minimum value(0)`
			// the global real-time index should also synchronize the same sliding distance
			m.index -= pageStart - 0
		}
	}
}

// forward triggers a fast jump action, if the pageIndex
// is invalid, keep it as it is
func (m *Model) forward(pageIndex string) {
	// the caller guarantees that pageIndex is an integer, and err is not processed here
	idx, _ := strconv.Atoi(pageIndex)
	idx--

	// pageIndex has exceeded the maximum index of the page, ignore
	if idx > m.pageMaxIndex {
		return
	}

	// calculate the distance moved to pageIndex
	l := idx - m.pageIndex
	// update the global real time index
	m.index += l
	// update the page real time index
	m.pageIndex = idx

}

// initData initialize the data model, set the default value and
// fix the wrong parameter settings during initialization
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
		m.HeaderFunc = func(_ Model, _ interface{}, _ int) string {
			return common.FontColor(DefaultHeader, ColorHeader)
		}
	}
	if m.Cursor == "" {
		m.Cursor = DefaultCursor
	}
	if m.CursorColor == "" {
		m.CursorColor = ColorCursor
	}
	if m.SelectedFunc == nil {
		m.SelectedFunc = func(m Model, obj interface{}, gdIndex int) string {
			return common.FontColor(fmt.Sprint(obj), ColorSelected)
		}
	}
	if m.UnSelectedFunc == nil {
		m.UnSelectedFunc = func(m Model, obj interface{}, gdIndex int) string {
			return common.FontColor(fmt.Sprint(obj), ColorUnSelected)
		}
	}
	if m.FooterFunc == nil {
		m.FooterFunc = func(_ Model, _ interface{}, _ int) string {
			return common.FontColor(DefaultFooter, ColorFooter)
		}
	}
	if m.FinishedFunc == nil {
		m.FinishedFunc = func(s interface{}) string {
			return common.FontColor(fmt.Sprintf(DefaultFinished, s), ColorFinished)
		}
	}
	m.init = true
}

// pageIndexInfo return the start and end positions of the slice of the
// page data area corresponding to the global data area
func (m *Model) pageIndexInfo() (start, end int) {
	// `Global real-time index` - `page real-time index` = `start index of page data area`
	start = m.index - m.pageIndex
	// `Page data area start index` + `single page size` = `page data area end index`
	end = start + m.PerPage
	return
}

// DefaultHeaderFuncWithAppend return the default HeaderFunc and append
// the given string to the next line of the default header
func DefaultHeaderFuncWithAppend(append string) func(m Model, obj interface{}, gdIndex int) string {
	return func(m Model, obj interface{}, gdIndex int) string {
		return common.FontColor(DefaultHeader+"\n"+append, ColorHeader)
	}
}

// DefaultSelectedFuncWithIndex return the default SelectedFunc and adds
// the serial number prefix of the given format
func DefaultSelectedFuncWithIndex(indexFormat string) func(m Model, obj interface{}, gdIndex int) string {
	return func(m Model, obj interface{}, gdIndex int) string {
		return common.FontColor(fmt.Sprintf(indexFormat+" %v", gdIndex+1, obj), ColorSelected)
	}
}

// DefaultUnSelectedFuncWithIndex return the default UnSelectedFunc and
// adds the serial number prefix of the given format
func DefaultUnSelectedFuncWithIndex(indexFormat string) func(m Model, obj interface{}, gdIndex int) string {
	return func(m Model, obj interface{}, gdIndex int) string {
		return common.FontColor(fmt.Sprintf(indexFormat+" %v", gdIndex+1, obj), ColorUnSelected)
	}
}

// Index return the global real time index
func (m Model) Index() int {
	return m.index
}

// PageIndex return the real time index of the page
func (m Model) PageIndex() int {
	return m.pageIndex
}

// PageData return the current page data area slice
func (m Model) PageData() []interface{} {
	return m.pageData
}

// Selected return the currently selected data
func (m Model) Selected() interface{} {
	return m.Data[m.index]
}

//// PageSelected return the currently selected data(same as the Selected func)
//func (m Model) PageSelected() interface{} {
//	return m.pageData[m.pageIndex]
//}

// Canceled determine whether the operation is cancelled
func (m Model) Canceled() bool {
	return m.canceled
}
