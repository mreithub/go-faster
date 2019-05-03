package internal

// KeyHTML -- dashboard's key info template
var KeyHTML = `
<html>
<head>
<title>{{.keyName}} :: go-faster key stats</title>
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.4.1/jquery.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/flot/0.8.3/jquery.flot.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/flot/0.8.3/jquery.flot.resize.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/flot/0.8.3/jquery.flot.time.js"></script>

<style>

</style>
</head>
<body>
<h2>go-faster key stats:
  {{range .keyPath}}
    <tt style="color: #aaa">{{.}} |</tt>
  {{end}}
  <tt>{{.keyName}}</tt>
</h2>

<a href="./">Back</a>

<button onclick="fetchData()">Reload</button>

<h3>Requests</h3>

<div>
  Ticker:
{{range .sortedTickers }}
  <a {{if ne $.ticker.Name .Name}}href="{{($.url.WithPath "key").WithParam "ticker" .Name}}"{{end}} title="last {{.Capacity}} snapshots (with interval {{.Interval}})">last {{.Duration}}</a>
{{end }}
</div>

<div id="chart" style="width: 100%; min-height: 300px;"></div>

</body>
<script>
function fetchData() {
  $.getJSON('{{.url.WithPath "key/history.json"}}', function(data) {
    document._hist = data;
    console.log('data: ', data);

    if (data.ts == null || data.ts.length == 0) {
      if ($('#chart .nodata').length == 0) {
        $('#chart').append('<div class="nodata">:: no data ::</div>');
      }
      return;
    } else $('#chart .nodata').remove();

    // format data the way flot expects it
    var counts = [], avgMsec = [];
    for (var i = 0; i < data.ts.length; i++) {
      counts.push([data.ts[i], data.counts[i]])
      avgMsec.push([data.ts[i], data.avgMsec[i]])
    }

    $.plot($("#chart"), [
        {
          data: counts,
          label: "requests",
          bars: {show: true, barWidth: 800, align: "center"},
        },
        {
          data: avgMsec,
          label: "average duration",
          yaxis: 2,
        },
      ], {
      xaxis: {
        mode: "time",
        timeBase: "milliseconds",
      },
      yaxes: [
        {min: 0},
        {
          min: 0,
          alignTicksWithAxis: 1,
          position: "right",
          tickFormatter: function(v, axis) {
            return v.toFixed(axis.tickDecimals) + "ms";
          },
        }
      ],
      /*legend: {
        backgroundColor: "transparent"
      }*/
    });

  })
}

fetchData();
</script>
</html>
`
