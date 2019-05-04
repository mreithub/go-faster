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

<h3>Histogram</h3>
<div id="histogram" style="width: 100%; min-height: 300px;"></div>

</body>
<script>
function fetchData() {
  $.getJSON('{{.url.WithPath "key/info.json"}}', function(data) {
    var req = data.requests;
    if (req.ts == null || req.ts.length == 0) {
      if ($('#chart .nodata').length == 0) {
        $('#chart').append('<div class="nodata">:: no data ::</div>');
      }
      return;
    } else $('#chart .nodata').remove();

    // format data the way flot expects it
    var counts = [], avgMsec = [];
    for (var i = 0; i < req.ts.length; i++) {
      counts.push([req.ts[i], req.counts[i]])
      avgMsec.push([req.ts[i], req.avgMsec[i]])
    }

    var histogram = [];
    for (var h of data.histogram) {
      histogram.push([Math.log2(h.ns), h.count])
    }

    $.plot($("#chart"), [
        {
          data: counts,
          label: "# of calls",
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
    });

    $.plot($("#histogram"), [
        {
          data: histogram,
          label: "",
          bars: {show: true, align: "center"},
        },
      ], {
      xaxis: {
        //mode: "time",
        //timeBase: "milliseconds",
        tickFormatter: function(v, axis) {
          var v = Math.pow(2, v)
          var units = ['ns', 'us', 'ms', 's']
          var unit = 0;
          for (var i = 0; i < units.length; i++) {
            if (v <= 1000) break;
            v /= 1000;
            unit++;
          }
          unit = units[unit];

          return v.toFixed(axis.tickDecimals) + "ms";
        },
      },
      yaxis: {
        min: 0
      },
    });
  })
}

fetchData();
</script>
</html>
`
