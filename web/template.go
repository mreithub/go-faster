package web

import "html/template"

var indexHTML = `
<html>
<head><title>Go-Faster statistics</title>
<style>
body {
  font-family: monospace;
}

td {
  text-align: right;
}
tr>td {
  text-align: initial;
}
</style>
</head>
<body>
<table>
  <thead><tr>
    <th>Name</th>
    <th title="number of currently running instances">Active</th>
    <th title="number of finished instances">Count</th>
    <th title="total time spent">Total Time</th>
    <th title="average time spent">Average Time</th>
  </tr></thead>
  <tbody>
    {{range .data}}
    <tr data-path="{{.Path}}">
      <td><span style="color: #aaa">
        {{range .Path}}{{.}}.{{end}}</span>{{.Name}}</td>
      <td>{{.Data.Active}}</td>
      <td>{{.Data.Count}}</td>
      <td data-raw="{{printf "%d" .Data.Duration}}">{{.Data.Duration}}</td>
      <td data-raw="{{printf "%d" .Data.AvgMsec}}">{{.Data.AvgMsec}}</td>
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
