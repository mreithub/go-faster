# GoFaster - Find bottlenecks in production Go Apps

[![Build Status](https://travis-ci.org/mreithub/faster.svg?branch=master)](https://travis-ci.org/mreithub/faster)


It can be used to:
- track execution time of individual functions/goroutines
- find bottlenecks in your code
- Check if your goroutines exit properly
- Track calls to your HTTP endpoints (and their execution time) - see below

To access the internal profiling data, use `GetSnapshot()`.
It'll ask the worker goroutine to create a deep copy of Faster's instance current state.

GoFaster's code is thread safe. It uses a messaging channel read by a single worker goroutine
which does the heavy lifting.  
Calls to `Ref()` and `Deref()` are asynchronous
(that asynchronousity doesn't affect time measurement though).  



## Getting started

Download the package, e.g.:

    go get github.com/mreithub/go-faster

Add the following snippet to each function (or goroutine) you want to track
(and replace 'foo' with your own key names).

```go
ref := faster.Ref("foo"); defer ref.Deref()
```

The above snippet uses `GoFaster` in singleton mode. But you can also create your
own `GoFaster` instances (and e.g. use different ones in different parts of your
application):

```go
g := faster.New()

// and then instead of the above snippet:
ref := g.Ref("foo"); defer ref.Deref()
```


At any point in time you can call `GetSnapshot()` to obtain a deep copy of the measurements.



### Scoped measurements

GoFaster not only supports independent `Faster` instances but also has a scope hierarchy (or tree structure
if you will).

With `faster.GetInstance(path ...string)` you can get a specific child of the global singleton instance.

An example use case would be seperate, possibly nested instances for different parts of your application
(e.g. `faster.GetInstance("http")` for HTTP endpoint handlers, `faster.GetInstance("dao", "psql")` for the PostgreSQL based DAO, ...)

You can see a simple example of GoFaster scopes in action in the *gorilla-mux* example below (or in the `examples/gorillamux/` directory)



## Example (excerpt from [webserver.go](examples/webserver/webserver.go)):

This example shows how to use GoFaster in your web applications.  
Here it tracks all web handler invocations.

Have a look at the usage documentation at [godoc.org][godoc].

```go
func indexHTML(w http.ResponseWriter, r *http.Request) {
	ref := faster.Ref("/")
	defer ref.Deref()

	w.Write([]byte(`<h1>Index</h1>
  <a href="/delayed.html">delayed.html</a><br />
  <a href="/faster.json">faster.json</a>`))
}

func delayedHTML(w http.ResponseWriter, r *http.Request) {
	ref := faster.Ref("/hello.html")
	defer ref.Deref()

	time.Sleep(200 * time.Millisecond)
	msg := fmt.Sprintf("The time is %s", time.Now().String())
	w.Write([]byte(msg))
}

func fasterJSON(w http.ResponseWriter, r *http.Request) {
	ref := faster.Ref("/faster.json")
	defer ref.Deref()

	data, _ := json.Marshal(faster.GetSnapshot().Data)

	w.Header().Add("Content-type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/", indexHTML)
	http.HandleFunc("/delayed.html", delayedHTML)
	http.HandleFunc("/faster.json", fasterJSON)

	http.ListenAndServe("localhost:1234", nil)
}
```

Run it with

    go run examples/webserver.go

and browse to http://localhost:1234/

After accessing each page a couple of times `/faster.json` should look something
like this:

```json
{
  "/": {
    "active": 0,
    "count": 6,
    "duration": 31131,
    "avgMsec": 0.0051885
  },
  "/delayed.html": {
    "active": 0,
    "count": 4,
    "duration": 811560843,
    "avgMsec": 202.89021
  },
  "/faster.json": {
    "active": 1,
    "count": 6,
    "duration": 443599,
    "avgMsec": 0.07393317
  }
}
```

- `active`: the number of currently active instances
- `count`: number of (finished) instances (doesn't include the `active` ones yet)
- `duration`: total time spent in that function (as time.Duration field)
- `avgMsec`: calculated average (`usec/(1000*total)`)



## Using [`gorilla-mux`][gorillamux]

If you're using [gorilla-mux][gorillamux], there's a simple way to
add GoFaster to your project:

(taken from the example in `examples/gorillamux/`)

```go
func trackRequests(router *mux.Router) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Try to find the matching HTTP route (we'll use that as GoFaster key)
    var match mux.RouteMatch
    if router.Match(r, &match) {
      path, _ := match.Route.GetPathTemplate()
      path = fmt.Sprintf("%s %s", r.Method, path)

      ref := faster.Ref(path)
      router.ServeHTTP(w, r)
      ref.Deref()
    } else {
      // No route found (-> 404 error)
      router.ServeHTTP(w, r)
    }
  })
}
```

and in your main function something like:

```go
var router = mux.NewRouter()
// add your routes here using router.HandleFunc() and the like
var addr = ":8080"
var handler = handlers.LoggingHandler(os.Stdout, trackRequests(router))
log.Fatal(http.ListenAndServe(addr, handler))
```

You'll get GoFaster data looking something like this:

```json
{
  "_children": {
    "app": {
      "data": {
        "processing": {
          "active": 0,
          "count": 11,
          "duration": 2232951708,
          "avgMsec": 202.9956
        }
      },
      "ts": "2017-06-02T21:50:11.717071564+02:00"
    },
    "http": {
      "data": {
        "GET /": {
          "active": 0,
          "count": 13,
          "duration": 193220,
          "avgMsec": 0.014863077
        },
        "GET /delayed.html": {
          "active": 0,
          "count": 11,
          "duration": 2233380060,
          "avgMsec": 203.03455
        },
        "GET /faster.json": {
          "active": 1,
          "count": 4,
          "duration": 2025613,
          "avgMsec": 0.50640327
        }
      },
      "ts": "2017-06-02T21:50:11.71706162+02:00"
    }
  },
  "ts": "2017-06-02T21:50:11.717049391+02:00"
}
```

Requests matched by the same gorilla-mux route will be grouped together.



## Performance impact

GoFaster aims to have as little impact on your application's performance as possible.

That's why all the processing is done asynchronously in a separate goroutine.

In a benchmark run on my laptop, this typical ref counter snippet takes around
a microsecond to run:

```go
r := faster.Ref(); defer r.Deref()
```

Interestingly, things are a lot faster if we don't use `defer`
as seen when running the `bench_test.go` benchmarks:

```
$ go test --run=XXX --bench=.
BenchmarkMeasureTime-4        	50000000	        33.9 ns/op
BenchmarkRefDeref-4           	 5000000	       339 ns/op
BenchmarkRefDerefDeferred-4   	 1000000	      1124 ns/op
BenchmarkGetSnapshot100-4     	  100000	     12367 ns/op
BenchmarkGetSnapshot1000-4    	   10000	    127117 ns/op
PASS
ok  	github.com/mreithub/faster	7.605s
```

- `BenchmarkMeasureTime()` measures the cost of calling time.Now() twice and calculating the nanoseconds between them
- `BenchmarkRefDeref()` calls `faster.Ref("hello").Deref()` directly (without using `defer`)
- `BenchmarkRefDerefDeferred()` uses `defer` (as in the snippet above)
- `BenchmarkGetSnapshot*()` measure the time it takes to take a snapshot of a GoFaster instance with 100 and 1000 entries (= different keys) respectively

[golang]: https://golang.org/
[godoc]: https://godoc.org/github.com/mreithub/faster
[gorillamux]: https://github.com/gorilla/mux
