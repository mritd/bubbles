package common

import (
	"github.com/muesli/termenv"
)

var term = termenv.ColorProfile()

// FontColor 对给定的字符串设置颜色并加粗字体
func FontColor(str, color string) string {
	return termenv.String(str).Foreground(term.Color(color)).Bold().String()
}

// GenSpaces 生成给定长度的空格
func GenSpaces(l int) string {
	var s string
	for i := 0; i < l; i++ {
		s += " "
	}
	return s
}
