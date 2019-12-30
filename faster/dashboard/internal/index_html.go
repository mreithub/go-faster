package internal

// IndexHTML -- dashboard index template
var IndexHTML = `
<html>
<head><title>go-faster dashboard</title>
<style>
body {
  font-family: monospace;
}

th, td { padding-left: 1em;}
tr:hover { background-color: rgba(192,224,255,.5);}

td { text-align: right; }
td:first-child { text-align: initial; }
</style>
</head>
<body>
<h1>go-faster dashboard</h1>


<h2>app info</h2>
<table><tbody>
<tr><th>hostname</th><td>{{.hostname}}</td></tr>
<tr><th>app uptime</th><td title="{{.startTS}}">{{.uptime}}</td></tr>
<tr><th>cpu</th><td>{{.cores}} cores</td></tr>
<tr><th>goroutines</th><td>{{.goroutines}}</td></tr>
</tbody></table>


<h2>stats</h2>
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
      <td data-raw="{{printf "%d" .Data.TotalTime}}" title="{{.Data.TotalTime}}">{{.PrettyTotal}}</td>
      <td data-raw="{{printf "%d" .Data.Average}}" title="{{.Data.Average}}">{{.PrettyAverage}}</td>
    </tr>
    {{end}}
  </tbody>
</table>
</body>
</html>
`
