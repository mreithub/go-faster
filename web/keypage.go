package web

import (
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/mreithub/go-faster/faster"
)

// KeyPage -- impl
type KeyPage struct {
	faster    *faster.Faster
	templates map[string]*template.Template
}

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

	var tickers = p.faster.ListTickers()
	var sortedTickers = p.sortHistoryByInterval(tickers)
	var selectedTicker *faster.History
	var data []*faster.Snapshot

	if len(sortedTickers) > 0 {
		selectedTicker = tickers[sortedTickers[0].Name] // TODO allow the user to pick another ticker
		data = selectedTicker.ForKey(key...).Relative()
	}

	var tpl = p.templates["key.html"]
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
