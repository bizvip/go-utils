package htm

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/tdewolff/minify/v2"
	minhtml "github.com/tdewolff/minify/v2/html"
)

// Compress 压缩 HTML 字符串。
//   - htmlSrc: 原始 HTML 字符串
//   - stripScriptStyle: 是否删除 <script> 与 <style> 标签
//
// 返回压缩后的 HTML 字符串。
func Compress(htmlSrc string, stripScriptStyle bool) (string, error) {
	// 初始化 minifier（如需高频调用，可放包级变量复用）
	m := minify.New()
	m.AddFunc("text/html", minhtml.Minify)

	// 1) 不过滤脚本/样式，直接压缩
	if !stripScriptStyle {
		minified, err := m.String("text/html", htmlSrc)
		return minified, err
	}

	// 2) 先用 goquery 删除 <script><style> 节点
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(htmlSrc))
	if err != nil {
		return "", err
	}
	doc.Find("script,style").Remove()

	cleanedHTML, err := doc.Html()
	if err != nil {
		return "", err
	}

	// 3) 压缩并返回
	minified, err := m.String("text/html", cleanedHTML)
	return minified, err
}
