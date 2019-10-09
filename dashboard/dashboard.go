package dashboard

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/mreithub/go-faster/faster"
)

// Dashboard -- implements go-faster's web dashboard
type Dashboard struct {
	faster    *faster.Faster
	mux       *http.ServeMux
	templates map[string]*template.Template
	keyPage   *keyPage

	AuthHandler http.Handler
}

// checkMethod -- returns true if the Request method's in the whitelist
func checkMethod(w http.ResponseWriter, r *http.Request, whitelist ...string) bool {
	for _, method := range whitelist {
		if r.Method == method {
			return true
		}
	}

	http.Error(w, "Method not allowed: "+r.Method, http.StatusMethodNotAllowed)
	return false
}

func (d *Dashboard) indexPage(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}
	ref := d.faster.Track("_faster", "index")
	defer ref.Done()

	var tpl = d.templates["index.html"]
	var data = flattenSnapshot(d.faster.TakeSnapshot())
	sortByPath(data)

	var hostname, _ = os.Hostname()
	var err = tpl.Execute(w, map[string]interface{}{
		"data":       data,
		"cores":      runtime.NumCPU(),
		"goroutines": runtime.NumGoroutine(),
		"hostname":   hostname,
		"uptime":     time.Now().Sub(d.faster.StartTS),
		"startTS":    d.faster.StartTS.Format(time.RFC3339),
	})

	if err != nil {
		log.Print("Error: failed to render go-faster index.html template: ", err.Error())
	}
}

func (d *Dashboard) snapshotJSON(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}
	ref := d.faster.Track("_faster", "snapshot.json")
	defer ref.Done()

	data, _ := json.MarshalIndent(d.faster.TakeSnapshot(), "", "  ")

	w.Header().Add("Content-type", "application/json")
	w.Write(data)
}

// ServeHTTP -- implements http.Handler
func (d *Dashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mux.ServeHTTP(w, r)
}

// New -- returns a HTTP Dashboard for the given Faster instance
func New(faster *faster.Faster) *Dashboard {
	var mux = http.NewServeMux()
	var templates, err = parseTemplates()
	if err != nil {
		panic(err) // this only happens if there are template parsing errors
	}
	var rc = Dashboard{
		faster:    faster,
		mux:       mux,
		templates: templates,
		keyPage: &keyPage{
			faster:    faster,
			templates: templates,
		},
	}

	mux.HandleFunc("/", rc.indexPage)
	mux.Handle("/key", rc.keyPage)
	mux.HandleFunc("/key/info.json", rc.keyPage.InfoJSON)
	mux.HandleFunc("/snapshot.json", rc.snapshotJSON)

	return &rc
}
