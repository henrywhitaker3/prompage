package views

import (
	"fmt"
	"html/template"
	"io"

	chartrender "github.com/go-echarts/go-echarts/v2/render"
)

var baseTpl = `
<div class="container h-full">
    <div class="item" id="{{ .ChartID }}" style="width:100%;min-height:100%;"></div>
</div>
{{- range .JSAssets.Values }}
   <script src="{{ . }}"></script>
{{- end }}
<script type="text/javascript">
    "use strict";
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}");
    let option_{{ .ChartID | safeJS }} = {{ .JSON }};
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
    {{- range .JSFunctions.Fns }}
    {{ . | safeJS }}
    {{- end }}
	new ResizeObserver(() => goecharts_{{ .ChartID | safeJS }}.resize()).observe(document.querySelector('#{{ .ChartID | safeJS }}'))
</script>
`

type chartRenderer struct {
	c      any
	before []func()
}

func newChartRenderer(c any, before ...func()) chartrender.Renderer {
	return &chartRenderer{
		c:      c,
		before: before,
	}
}

func (r *chartRenderer) Render(w io.Writer) error {
	const tplName = "chart"
	for _, fn := range r.before {
		fn()
	}

	tpl := template.
		Must(template.New(tplName).
			Funcs(template.FuncMap{
				"safeJS": func(s interface{}) template.JS {
					return template.JS(fmt.Sprint(s))
				},
			}).
			Parse(baseTpl),
		)

	err := tpl.ExecuteTemplate(w, tplName, r.c)
	return err
}
