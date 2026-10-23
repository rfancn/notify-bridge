package zentao

import "testing"

func TestHTMLToMarkdown(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"空内容", "  ", ""},
		{"纯文本", "纯文本内容", "纯文本内容"},
		{"段落与换行", "<p>第一行<br/>第二行</p>", "第一行\n\n第二行"},
		{"加粗与实体", "<p><strong>重点</strong>内容&nbsp;说明</p>", "**重点** 内容 说明"},
		{"斜体", "<div><em>斜体</em></div>", "_斜体_"},
		{"链接", `<a href="http://example.com?a=1&amp;b=2">链接</a>`, "[链接](http://example.com?a=1&b=2)"},
		{"列表", "<ul><li>一</li><li>二</li></ul>", "- 一\n- 二"},
		{"段落分隔", "<p>a</p><p>b</p>", "a\n\nb"},
		{"代码", "<pre><code>x := 1</code></pre>", "```\nx := 1\n```"},
		{"图片", `<img src="http://example.com/a.png" alt="截图"/>`, "![截图](http://example.com/a.png)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := htmlToMarkdown(tt.src); got != tt.want {
				t.Errorf("htmlToMarkdown(%q) = %q, want %q", tt.src, got, tt.want)
			}
		})
	}
}
