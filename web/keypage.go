package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/mreithub/go-faster/faster"
)

type KeyPage struct {
	faster    *faster.Faster
	templates map[string]*template.Template
}

// HistoryJSON -- implements GET key/history.json
func (p *KeyPage) HistoryJSON(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}

	var key = r.URL.Query()["k"]
	if len(key) == 0 {
		http.Error(w, "missing key parameter(s): 'k'", http.StatusBadRequest)
		return
	}

	ref := p.faster.Track("_faster", "key", "history.json")
	defer ref.Done()

	var tickers = p.faster.ListTickers()
	var sortedTickers = p.sortHistoryByInterval(tickers)
	var selectedTicker *faster.History

	type Response struct {
		TS      []int64 `json:"ts"`
		Counts  []int64 `json:"counts"`
		AvgMsec []int64 `json:"avgMsec"`
	}
	var response Response

	if len(sortedTickers) > 0 {
		selectedTicker = tickers[sortedTickers[0].Name] // TODO allow the user to pick another ticker
		for _, snap := range selectedTicker.ForKey(key...).Relative() {
			response.TS = append(response.TS, snap.Ts.UnixNano()/int64(time.Millisecond))
			response.Counts = append(response.Counts, snap.Count)
			response.AvgMsec = append(response.AvgMsec, int64(snap.Average/time.Millisecond))
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Print("Error: failed to encode history.json: ", err)
	}
}

// ServeHTTP -- implements GET key
func (p *KeyPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}

	var key = r.URL.Query()["k"]
	if len(key) == 0 {
		http.Error(w, "missing key parameter(s): 'k'", http.StatusBadRequest)
		return
	}

	ref := p.faster.Track("_faster", "key")
	defer ref.Done()

	var tpl = p.templates["key.html"]
	var err = tpl.Execute(w, map[string]interface{}{
		"keyPath":  key[:len(key)-1],
		"keyName":  key[len(key)-1],
		"rawQuery": r.URL.RawQuery,
	})

	if err != nil {
		log.Print("Error: failed to render go-faster key.html template: ", err.Error())
	}

}

// sorts the given History tickers by their interval - lowest first
func (p *KeyPage) sortHistoryByInterval(data map[string]*faster.History) []*faster.History {
	if data == nil || len(data) == 0 {
		return nil
	}
	var rc = make([]*faster.History, 0, len(data))

	for _, h := range data {
		rc = append(rc, h)
	}

	sort.Slice(rc, func(i int, j int) bool {
		return rc[i].Interval() < rc[j].Interval()
	})
	return rc
}
