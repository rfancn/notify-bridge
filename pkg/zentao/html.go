package zentao

import (
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
)

// htmlToMarkdown 将禅道富文本备注转换为 Mattermost 可渲染的 Markdown。
func htmlToMarkdown(src string) string {
	if strings.TrimSpace(src) == "" {
		return ""
	}

	converter := htmltomarkdown.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())

	markdown, err := converter.ConvertString(src)
	if err != nil {
		// 转换失败时退化为去标签的纯文本
		return normalizeMarkdown(stripTags(src))
	}
	return normalizeMarkdown(markdown)
}

// stripTags 去掉 HTML 标签，仅在转换失败时兜底。
func stripTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// normalizeMarkdown 统一换行与行尾空白，保留缩进以维持列表层级和代码块。
func normalizeMarkdown(s string) string {
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
