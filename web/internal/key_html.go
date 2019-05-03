package internal

// KeyHTML -- dashboard's key info template
var KeyHTML = `
<html>
<head>
<title>{{.keyName}} :: go-faster key stats</title>
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/flot/0.8.3/jquery.flot.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/flot/0.8.3/jquery.flot.resize.min.js"></script>
<style>
</style>
</head>
<body>
<h2>go-faster key stats: {{.keyName}}

<h3>Tickers</h3>
{{range .tickers}}
  <!-- TODO these links obviously don't work -->
  <a href="?ticker={{.Name}}">{{ .Duration }}</a>
{{end}}

<h3>Requests</h3>
{{range .data}}
  <div>{{.}}</div>
{{end}}
</body>
</html>
`
