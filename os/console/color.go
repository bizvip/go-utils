package console

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

// C 控制台输出结构体
type C struct {
	writer io.Writer
}

// Console 返回默认的控制台实例（输出到 stdout）
func Console() *C {
	return &C{writer: os.Stdout}
}

// ConsoleErr 返回错误输出的控制台实例（输出到 stderr）
func ConsoleErr() *C {
	return &C{writer: os.Stderr}
}

// WithWriter 设置自定义的输出流
func (c *C) WithWriter(w io.Writer) *C {
	return &C{writer: w}
}

// ========== 基础颜色方法 ==========

// Black 黑色文本
func (c *C) Black(txt string) {
	c.print(color.FgHiBlack, txt)
}

// BlackBold 黑色加粗文本（白底黑字）
func (c *C) BlackBold(txt string) {
	c.printWithStyle(txt, color.FgHiWhite, color.BgBlack)
}

// Red 红色文本
func (c *C) Red(txt string) {
	c.print(color.FgHiRed, txt)
}

// RedBold 红色加粗文本
func (c *C) RedBold(txt string) {
	c.printWithStyle(txt, color.FgHiRed, color.Bold)
}

// Green 绿色文本
func (c *C) Green(txt string) {
	c.print(color.FgHiGreen, txt)
}

// GreenBold 绿色加粗文本
func (c *C) GreenBold(txt string) {
	c.printWithStyle(txt, color.FgHiGreen, color.Bold)
}

// Yellow 黄色文本
func (c *C) Yellow(txt string) {
	c.print(color.FgHiYellow, txt)
}

// YellowBold 黄色加粗文本
func (c *C) YellowBold(txt string) {
	c.printWithStyle(txt, color.FgHiYellow, color.Bold)
}

// Blue 蓝色文本
func (c *C) Blue(txt string) {
	c.print(color.FgHiBlue, txt)
}

// BlueBold 蓝色加粗文本
func (c *C) BlueBold(txt string) {
	c.printWithStyle(txt, color.FgHiBlue, color.Bold)
}

// Magenta 品红色文本
func (c *C) Magenta(txt string) {
	c.print(color.FgHiMagenta, txt)
}

// MagentaBold 品红色加粗文本
func (c *C) MagentaBold(txt string) {
	c.printWithStyle(txt, color.FgHiMagenta, color.Bold)
}

// Cyan 青色文本
func (c *C) Cyan(txt string) {
	c.print(color.FgHiCyan, txt)
}

// CyanBold 青色加粗文本
func (c *C) CyanBold(txt string) {
	c.printWithStyle(txt, color.FgHiCyan, color.Bold)
}

// White 白色文本
func (c *C) White(txt string) {
	c.print(color.FgHiWhite, txt)
}

// WhiteBold 白色加粗文本
func (c *C) WhiteBold(txt string) {
	c.printWithStyle(txt, color.FgHiWhite, color.Bold)
}

// Gray 灰色文本
func (c *C) Gray(txt string) {
	c.print(color.FgWhite, txt)
}

// GrayBold 灰色加粗文本
func (c *C) GrayBold(txt string) {
	c.printWithStyle(txt, color.FgWhite, color.Bold)
}

// ========== 格式化输出方法 ==========

// Printf 格式化输出（无颜色）
func (c *C) Printf(format string, a ...interface{}) {
	fmt.Fprintf(c.writer, format, a...)
}

// Redf 红色格式化输出
func (c *C) Redf(format string, a ...interface{}) {
	c.Red(fmt.Sprintf(format, a...))
}

// Greenf 绿色格式化输出
func (c *C) Greenf(format string, a ...interface{}) {
	c.Green(fmt.Sprintf(format, a...))
}

// Yellowf 黄色格式化输出
func (c *C) Yellowf(format string, a ...interface{}) {
	c.Yellow(fmt.Sprintf(format, a...))
}

// Bluef 蓝色格式化输出
func (c *C) Bluef(format string, a ...interface{}) {
	c.Blue(fmt.Sprintf(format, a...))
}

// Cyanf 青色格式化输出
func (c *C) Cyanf(format string, a ...interface{}) {
	c.Cyan(fmt.Sprintf(format, a...))
}

// Whitef 白色格式化输出
func (c *C) Whitef(format string, a ...interface{}) {
	c.White(fmt.Sprintf(format, a...))
}

// Grayf 灰色格式化输出
func (c *C) Grayf(format string, a ...interface{}) {
	c.Gray(fmt.Sprintf(format, a...))
}

// Magentaf 品红色格式化输出
func (c *C) Magentaf(format string, a ...interface{}) {
	c.Magenta(fmt.Sprintf(format, a...))
}

// ========== 特殊样式方法 ==========

// Success 成功消息（绿色背景）
func (c *C) Success(txt string) {
	c.printWithStyle(" ✓ "+txt+" ", color.FgHiWhite, color.BgGreen)
}

// Error 错误消息（红色背景）
func (c *C) Error(txt string) {
	c.printWithStyle(" ✗ "+txt+" ", color.FgHiWhite, color.BgRed)
}

// Warning 警告消息（黄色背景）
func (c *C) Warning(txt string) {
	c.printWithStyle(" ⚠ "+txt+" ", color.FgHiBlack, color.BgYellow)
}

// Info 信息消息（蓝色背景）
func (c *C) Info(txt string) {
	c.printWithStyle(" ℹ "+txt+" ", color.FgHiWhite, color.BgBlue)
}

// Underline 下划线文本
func (c *C) Underline(txt string) {
	c.printWithStyle(txt, color.Underline)
}

// Italic 斜体文本
func (c *C) Italic(txt string) {
	c.printWithStyle(txt, color.Italic)
}

// ========== 进度和装饰方法 ==========

// Line 打印分隔线
func (c *C) Line(width int) {
	c.Gray(strings.Repeat("─", width) + "\n")
}

// DoubleLine 打印双分隔线
func (c *C) DoubleLine(width int) {
	c.Gray(strings.Repeat("═", width) + "\n")
}

// DashedLine 打印虚线分隔线
func (c *C) DashedLine(width int) {
	c.Gray(strings.Repeat("- ", width/2) + "\n")
}

// Box 打印带边框的文本
func (c *C) Box(txt string, color color.Attribute) {
	width := len(txt) + 4
	topBottom := strings.Repeat("─", width-2)

	c.printWithStyle("┌"+topBottom+"┐\n", color)
	c.printWithStyle("│ "+txt+" │\n", color)
	c.printWithStyle("└"+topBottom+"┘\n", color)
}

// Title 打印标题（带装饰）
func (c *C) Title(txt string) {
	c.WhiteBold("\n╔═══ " + txt + " ═══╗\n")
}

// Section 打印章节标题
func (c *C) Section(txt string) {
	c.CyanBold("\n▶ " + txt + "\n")
	c.Line(50)
}

// ========== 列表和表格方法 ==========

// List 打印列表项
func (c *C) List(items []string) {
	for _, item := range items {
		c.White("  • " + item + "\n")
	}
}

// NumberedList 打印编号列表
func (c *C) NumberedList(items []string) {
	for i, item := range items {
		c.Whitef("  %d. %s\n", i+1, item)
	}
}

// KeyValue 打印键值对
func (c *C) KeyValue(key, value string, keyWidth int) {
	format := fmt.Sprintf("  %%-%ds : %%s\n", keyWidth)
	c.White(fmt.Sprintf(format, key, ""))
	c.Gray(value + "\n")
}

// ========== 动画效果方法 ==========

// Spinner 显示加载动画（需要在goroutine中使用）
func (c *C) Spinner(message string, done <-chan bool) {
	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Fprint(c.writer, "\r\033[K") // 清除行
			return
		default:
			c.Cyanf("\r%s %s", spinChars[i%len(spinChars)], message)
			i++
			// 这里应该导入 time 包并添加延迟
			// time.Sleep(100 * time.Millisecond)
		}
	}
}

// Progress 显示进度条
func (c *C) Progress(current, total int, width int) {
	percent := float64(current) / float64(total)
	filled := int(percent * float64(width))
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	c.Greenf("\r[%s] %.0f%%", bar, percent*100)

	if current >= total {
		fmt.Fprintln(c.writer) // 完成时换行
	}
}

// ========== 内部辅助方法 ==========

// print 基础打印方法
func (c *C) print(attr color.Attribute, txt string) {
	col := color.New(attr).SprintFunc()
	fmt.Fprint(c.writer, col(txt))
}

// printWithStyle 带多个样式的打印方法
func (c *C) printWithStyle(txt string, attrs ...color.Attribute) {
	col := color.New(attrs...).SprintFunc()
	fmt.Fprint(c.writer, col(txt))
}

// Clear 清屏
func (c *C) Clear() {
	fmt.Fprint(c.writer, "\033[2J\033[H")
}

// NewLine 换行
func (c *C) NewLine() {
	fmt.Fprintln(c.writer)
}

// ========== 链式调用支持 ==========

// Print 打印文本（支持链式调用）
func (c *C) Print(txt string) *C {
	fmt.Fprint(c.writer, txt)
	return c
}

// Println 打印文本并换行
func (c *C) Println(txt string) *C {
	fmt.Fprintln(c.writer, txt)
	return c
}
