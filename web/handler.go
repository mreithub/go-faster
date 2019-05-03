package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

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

func (h *WebHandler) historyJSON(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}

	var data = map[string]interface{}{}
	for name, history := range h.faster.ListTickers() {
		data[name] = history.Len()
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Print("Failed to encode JSON response: ", err)
	}
}

func (h *WebHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}
	ref := h.faster.Track("_faster", "index")
	defer ref.Done()

	var tpl = h.templates["index.html"]
	var data = flattenSnapshot(faster.GetSnapshot())
	sortByPath(data)
	var err = tpl.Execute(w, map[string]interface{}{
		"data": data,
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

	data, _ := json.MarshalIndent(faster.GetSnapshot(), "", "  ")

	w.Header().Add("Content-type", "application/json")
	w.Write(data)
}

func (h *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// NewHandler -- returns a http handler for the given GoFaster instance
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

	mux.HandleFunc("/", rc.indexHandler)
	mux.Handle("/key", rc.keyPage)
	mux.HandleFunc("/key/history.json", rc.keyPage.HistoryJSON)
	mux.HandleFunc("/snapshot.json", rc.snapshotJSON)
	mux.HandleFunc("/history.json", rc.historyJSON)

	return &rc
}
