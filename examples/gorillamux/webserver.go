package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mreithub/go-faster/faster"
	"github.com/mreithub/go-faster/web"
)

func basicAuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer faster.TrackFn().Done()
		if user, pw, ok := r.BasicAuth(); ok && user == "admin" && pw == "hackme" {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Add("WWW-Authenticate", "Basic realm=\"stats\"")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}

func indexHTML(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<h1>Index</h1>
  <a href="/delayed.html">delayed.html</a><br />
  <a href="/_faster/">go-faster dashboard</a> (user: <tt>admin</tt> password: <tt>hackme</tt>)`))
}

func delayedHTML(w http.ResponseWriter, r *http.Request) {
	foo := processStuff(r.RemoteAddr)
	result := <-foo
	msg := fmt.Sprintf("Incoming message: %s", result)

	w.Header().Set("content-type", "text/html")
	w.Write([]byte(msg))
}

func processStuff(name string) chan string {
	rc := make(chan string)

	go func() {
		// since processing takes some time, we'll add a separate GoFaster instance here (this time in the "app" scope)
		defer faster.TrackFn().Done()

		var delay = time.Duration(rand.Intn(1700)) * time.Millisecond
		time.Sleep(delay)
		rc <- fmt.Sprintf("Hello %s, processing your request took <tt>%v</tt><br />\n<a href=\"javascript:history.back()\">back</a>", name, delay)
	}()

	return rc
}

func fasterJSON(w http.ResponseWriter, r *http.Request) {
	data, _ := json.MarshalIndent(faster.GetSnapshot(), "", "  ")

	w.Header().Add("Content-type", "application/json")
	w.Write(data)
}

func trackRequests(router *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to find the matching HTTP route (we'll use that as GoFaster key)
		var match mux.RouteMatch
		if router.Match(r, &match) {
			path, _ := match.Route.GetPathTemplate()
			path = fmt.Sprintf("%s %s", r.Method, path)

			ref := faster.Track("http", path)
			router.ServeHTTP(w, r)
			ref.Done()
		} else {
			// No route found (-> 404 error)
			router.ServeHTTP(w, r)
		}
	})
}

func main() {
	// setup http mux and loggin
	var r = mux.NewRouter()
	r.HandleFunc("/", indexHTML)
	r.HandleFunc("/delayed.html", delayedHTML)

	// add simple HTTP basic auth to the go-faster stats page (as it might expose sensitive info)
	var s = r.PathPrefix("/_faster").Subrouter()
	s.Use(basicAuthMW)
	s.NewRoute().Handler(http.StripPrefix("/_faster", web.NewHandler(faster.Singleton)))
	var handler = handlers.LoggingHandler(os.Stdout, trackRequests(r))

	// set up periodic go-faster snapshots
	faster.SetTicker("1sec", 1*time.Second, 120) // 2min
	faster.SetTicker("1min", 1*time.Minute, 60)  // 1h

	// start web server
	var addr = "localhost:1234"
	log.Printf("starting web server at '%s'", addr)
	log.Printf(" - go to 'http://%s/_faster/' for the dashboard", addr)

	log.Fatal(http.ListenAndServe(addr, handler))
}
