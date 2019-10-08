package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mreithub/go-faster/dashboard"
	"github.com/mreithub/go-faster/faster"
)

func indexHTML(w http.ResponseWriter, r *http.Request) {
	ref := faster.Track("http", r.Method+" /")
	defer ref.Done()

	w.Write([]byte(`<h1>Index</h1>
  <a href="/delayed.html">delayed.html</a><br />
  <a href="/_faster/">go-faster dashboard</a>`))
}

func delayedHTML(w http.ResponseWriter, r *http.Request) {
	ref := faster.Track("http", r.Method+" /delayed.html")
	defer ref.Done()

	time.Sleep(200 * time.Millisecond)
	msg := fmt.Sprintf("The time is %s", time.Now().String())
	w.Write([]byte(msg))
}

func main() {
	http.HandleFunc("/", indexHTML)
	http.HandleFunc("/delayed.html", delayedHTML)
	http.Handle("/_faster/", http.StripPrefix("/_faster", dashboard.NewHandler(faster.Singleton)))

	var addr = "localhost:1234"
	log.Printf("starting web server at '%s'", addr)
	log.Printf(" - go to 'http://%s/_faster/' for the dashboard", addr)
	http.ListenAndServe(addr, nil)
}
