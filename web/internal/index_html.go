package internal

// IndexHTML -- dashboard index template
var IndexHTML = `
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
