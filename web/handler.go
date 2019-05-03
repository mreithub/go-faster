package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/mreithub/go-faster/faster"
)

// WebHandler -- implements go-faster's web dashboard
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

func (h *WebHandler) historyJSON(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, "GET") {
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

func (h *WebHandler) keyHandler(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, "GET") {
		return
	}

	var key = r.URL.Query()["k"]
	if len(key) == 0 {
		http.Error(w, "missing key parameter(s): 'k'", http.StatusBadRequest)
		return
	}

	ref := h.faster.Track("_faster", "key")
	defer ref.Done()

	var tickers = h.faster.ListTickers()
	var sortedTickers = sortHistoryByInterval(tickers)
	var selectedTicker *faster.History
	var data []*faster.Snapshot

	if len(sortedTickers) > 0 {
		selectedTicker = tickers[sortedTickers[0].Name] // TODO allow the user to pick another ticker
		data = diff(selectedTicker.ForKey(key...))
	}

	var tpl = h.templates["key.html"]
	var err = tpl.Execute(w, map[string]interface{}{
		"key":     key,
		"keyName": key[len(key)-1],
		"tickers": sortedTickers,
		"data":    data,
	})

	if err != nil {
		log.Print("Error: failed to render go-faster key.html template: ", err.Error())
	}
}

func (h *WebHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, "GET") {
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
	if !h.checkMethod(w, r, "GET") {
		return
	}
	ref := h.faster.Track("_faster", "snapshot.json")
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
	var templates, err = parseTemplates(prefix)
	if err != nil {
		panic(err) // this only happens if there are template parsing errors
	}
	var rc = WebHandler{
		faster:    faster,
		prefix:    prefix,
		mux:       mux,
		templates: templates,
	}

	mux.HandleFunc(prefix, rc.indexHandler)
	mux.HandleFunc(path.Join(prefix, "key"), rc.keyHandler)
	mux.HandleFunc(path.Join(prefix, "snapshot.json"), rc.snapshotJSON)
	mux.HandleFunc(path.Join(prefix, "history.json"), rc.historyJSON)

	return &rc
}
