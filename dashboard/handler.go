package dashboard

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/mreithub/go-faster/faster"
)

// WebHandler -- implements go-faster's web dashboard
type WebHandler struct {
	faster    *faster.Faster
	mux       *http.ServeMux
	templates map[string]*template.Template
	keyPage   *KeyPage

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

func (h *WebHandler) indexPage(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}
	ref := h.faster.Track("_faster", "index")
	defer ref.Done()

	var tpl = h.templates["index.html"]
	var data = flattenSnapshot(faster.TakeSnapshot())
	sortByPath(data)
	var err = tpl.Execute(w, map[string]interface{}{
		"data":       data,
		"cores":      runtime.NumCPU(),
		"goroutines": runtime.NumGoroutine(),
		"uptime":     time.Now().Sub(h.faster.StartTS),
		"startTS":    h.faster.StartTS.Format(time.RFC3339),
	})

	if err != nil {
		log.Print("Error: failed to render go-faster index.html template: ", err.Error())
	}
}

func (h *WebHandler) snapshotJSON(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}
	ref := h.faster.Track("_faster", "snapshot.json")
	defer ref.Done()

	data, _ := json.MarshalIndent(faster.TakeSnapshot(), "", "  ")

	w.Header().Add("Content-type", "application/json")
	w.Write(data)
}

func (h *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// NewHandler -- returns a http handler for the given Faster instance
func NewHandler(faster *faster.Faster) http.Handler {
	var mux = http.NewServeMux()
	var templates, err = parseTemplates()
	if err != nil {
		panic(err) // this only happens if there are template parsing errors
	}
	var rc = WebHandler{
		faster:    faster,
		mux:       mux,
		templates: templates,
		keyPage: &KeyPage{
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
