package zentao

import "testing"

func TestHTMLToPlainText(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"空内容", "  ", ""},
		{"纯文本", "纯文本内容", "纯文本内容"},
		{"段落与换行", "<p>第一行<br/>第二行</p>", "第一行\n第二行"},
		{"加粗与实体", "<p><strong>重点</strong>内容&nbsp;说明</p>", "重点内容 说明"},
		{"斜体", "<div><em>斜体</em></div>", "斜体"},
		{"链接", `<a href="http://example.com?a=1&amp;b=2">链接</a>`, "链接"},
		{"列表", "<ul><li>一</li><li>二</li></ul>", "一\n二"},
		{"段落分隔", "<p>a</p><p>b</p>", "a\nb"},
		{"代码", "<pre><code>x := 1</code></pre>", "x := 1"},
		{"图片", `<img src="http://example.com/a.png" alt="截图"/>`, ""},
		{"内联样式", `<p><span style="color:red">红字</span></p>`, "红字"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := htmlToPlainText(tt.src); got != tt.want {
				t.Errorf("htmlToPlainText(%q) = %q, want %q", tt.src, got, tt.want)
			}
		})
	}
}
