package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mreithub/go-faster/faster"
)

func indexHTML(w http.ResponseWriter, r *http.Request) {
	ref := faster.Ref("/")
	defer ref.Deref()

	w.Write([]byte(`<h1>Index</h1>
  <a href="/delayed.html">delayed.html</a><br />
  <a href="/faster.json">faster.json</a>`))
}

func delayedHTML(w http.ResponseWriter, r *http.Request) {
	ref := faster.Ref("/delayed.html")
	defer ref.Deref()

	time.Sleep(200 * time.Millisecond)
	msg := fmt.Sprintf("The time is %s", time.Now().String())
	w.Write([]byte(msg))
}

func fasterJSON(w http.ResponseWriter, r *http.Request) {
	ref := faster.Ref("/faster.json")
	defer ref.Deref()

	data, _ := json.MarshalIndent(faster.GetSnapshot().Data, "", "  ")

	w.Header().Add("Content-type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/", indexHTML)
	http.HandleFunc("/delayed.html", delayedHTML)
	http.HandleFunc("/faster.json", fasterJSON)

	http.ListenAndServe("localhost:1234", nil)
}
