package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/mreithub/go-faster/faster"
)

type WebHandler struct {
	faster    *faster.Faster
	mux       *http.ServeMux
	prefix    string
	templates map[string]*template.Template

	// TODO add auth callback
}

// checkMethod -- returns true if the Request method's in the whitelist
func (h *WebHandler) checkMethod(w http.ResponseWriter, r *http.Request, whitelist ...string) bool {
	for _, method := range whitelist {
		if r.Method == method {
			return true
		}
	}

	http.Error(w, "Method not allowed: "+r.Method, http.StatusMethodNotAllowed)
	return false
}

func (h *WebHandler) indexHTML(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, "GET") {
		return
	}
	ref := h.faster.Track("http", "_faster", "GET index.html")
	defer ref.Done()

	var tpl = h.templates["index.html"]
	var data = flattenSnapshot(faster.GetSnapshot())
	sortByPath(data)
	var err = tpl.Execute(w, map[string]interface{}{
		"data": data,
	})

	if err != nil {
		log.Print("Error: failed to render go-faster template: ", err.Error())
	}
}

func (h *WebHandler) snapshotJSON(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, "GET") {
		return
	}
	ref := h.faster.Track("http", "_faster", "GET snapshot.json")
	defer ref.Done()

	data, _ := json.MarshalIndent(faster.GetSnapshot(), "", "  ")

	w.Header().Add("Content-type", "application/json")
	w.Write(data)
}

func (h *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO add auth here
	h.mux.ServeHTTP(w, r)
}

// NewHandler -- returns a http handler for the given GoFaster instance
func NewHandler(prefix string, faster *faster.Faster) http.Handler {
	var mux = http.NewServeMux()
	var templates, err = parseTemplates()
	if err != nil {
		panic(err) // this only happens if there are template parsing errors
	}
	var rc = WebHandler{
		faster:    faster,
		prefix:    prefix,
		mux:       mux,
		templates: templates,
	}

	mux.HandleFunc(prefix, rc.indexHTML)
	mux.HandleFunc(path.Join(prefix, "index.html"), rc.indexHTML)
	mux.HandleFunc(path.Join(prefix, "snapshot.json"), rc.snapshotJSON)

	return &rc
}
