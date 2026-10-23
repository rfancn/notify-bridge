package zentao

import (
	"html"
	"regexp"
	"strings"
)

var (
	// 块级标签与 <br> 转为换行，避免各段文字粘连。
	blockTagRe = regexp.MustCompile(`(?i)<\s*(?:br|/p|/div|/li|/tr|/h[1-6]|/blockquote|/pre)\s*/?\s*>`)
	htmlTagRe  = regexp.MustCompile(`<[^>]*>`)
)

// htmlToPlainText 将禅道富文本备注转换为纯文本，去除所有 HTML 标签和样式。
func htmlToPlainText(src string) string {
	if strings.TrimSpace(src) == "" {
		return ""
	}

	text := blockTagRe.ReplaceAllString(src, "\n")
	text = htmlTagRe.ReplaceAllString(text, "")
	text = html.UnescapeString(text)
	return normalizeText(text)
}

// normalizeText 统一换行与行尾空白，并合并多余空行。
func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	s = strings.Join(lines, "\n")

	for strings.Contains(s, "\n\n\n") {
		s = strings.ReplaceAll(s, "\n\n\n", "\n\n")
	}
	return strings.TrimSpace(s)
}
