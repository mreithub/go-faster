package web

import "html/template"

var indexHTML = `
<html>
<head><title>Go-Faster statistics</title>
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
      <td>{{range .Path}}&nbsp;&nbsp;{{end}}{{.Name}}</td>
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

func parseTemplates() (map[string]*template.Template, error) {
	var tpls = map[string]string{
		"index.html": indexHTML,
	}
	var rc = map[string]*template.Template{}
	var err error

	for name, html := range tpls {
		var tpl *template.Template
		if tpl, err = template.New(name).Parse(html); err != nil {
			return nil, err
		}
		rc[name] = tpl
	}

	return rc, nil
}
