package web

import (
	"html/template"
	"net/url"
	"path"
)

var indexHTML = `
<html>
<head><title>go-faster stats</title>
<style>
body {
  font-family: monospace;
}

th { padding-left: 1em;}

td { text-align: right; }
td:first-child { text-align: initial; }
</style>
</head>
<body>
<table>
  <thead><tr>
    <th>Name</th>
    <th title="number of currently running instances">active</th>
    <th title="number of finished instances">count</th>
    <th title="total time spent">total ms</th>
    <th title="average time spent">average ms</th>
  </tr></thead>
  <tbody>
    {{range .data}}
    <tr data-path="{{.JSONPath}}">
      <td>{{range .Path}}&nbsp;&nbsp;{{end -}}
        {{if gt .Data.Count 0 }}
          <a href="{{keyLink .Key}}">{{.Name}}</a>
        {{else}}
          {{.Name}}
        {{end}}
      </td>
      <td>{{or .Data.Active ""}}</td>
      <td>{{or .Data.Count ""}}</td>
      <td data-raw="{{printf "%d" .Data.Duration}}" title="{{.Data.Duration}}">{{.PrettyTotal}}</td>
      <td data-raw="{{printf "%d" .Data.Average}}" title="{{.Data.Average}}">{{.PrettyAverage}}</td>
    </tr>
    {{end}}
  </tbody>
</table>
</body>
</html>
`

var keyHTML = `
<html>
<head><title>{{.keyName}} - go-faster key stats</title>
<style>
</style>
</head>
<body>
</body>
</html>
`

func parseTemplates(prefix string) (map[string]*template.Template, error) {
	var tpls = map[string]string{
		"index.html": indexHTML,
		"key.html":   keyHTML,
	}
	var rc = map[string]*template.Template{}
	var err error

	var funcs = map[string]interface{}{
		"keyLink": func(key []string) string {
			var query = url.Values{
				"k": key,
			}
			var rc = url.URL{
				Path:     path.Join(prefix, "key"),
				RawQuery: query.Encode(),
			}
			return rc.String()
		},
	}

	for name, html := range tpls {
		var tpl *template.Template
		if tpl, err = template.New(name).Funcs(funcs).Parse(html); err != nil {
			return nil, err
		}
		rc[name] = tpl
	}

	return rc, nil
}
